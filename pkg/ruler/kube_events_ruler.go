package ruler

import (
	"context"
	"fmt"
	"sync"

	"github.com/kubesphere/kube-events/pkg/config"
	"github.com/kubesphere/kube-events/pkg/ruler/sinks/alert"
	"github.com/kubesphere/kube-events/pkg/ruler/sinks/notification"
	"github.com/kubesphere/kube-events/pkg/ruler/types"
	"github.com/panjf2000/ants/v2"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
)

type KubeEventsRuler struct {
	kcfg *rest.Config

	rulerConds      *rulerConds
	rulerCondsMutex sync.Mutex

	taskPool *ants.Pool

	evtQueue      workqueue.RateLimitingInterface
	notificaQueue workqueue.RateLimitingInterface
	alertQueue    workqueue.RateLimitingInterface
}

func (r *KubeEventsRuler) ReloadConfig(c *config.RulerConfig) error {
	r.rulerCondsMutex.Lock()
	defer r.rulerCondsMutex.Unlock()
	if c == nil {
		return nil
	}

	var notificaSinkers []types.NotificationSinker
	var alertSinkers []types.AlertSinker
	if c.Sinks != nil {
		for _, w := range c.Sinks.Webhooks {
			url := w.Url
			if url == "" && w.Service != nil {
				url = fmt.Sprintf("http://%s.%s.svc:%d/%s",
					w.Service.Name, w.Service.Namespace, *w.Service.Port, w.Service.Path)
			}
			if url == "" {
				continue
			}
			switch w.Type {
			case config.RulerSinkTypeNotification:
				notificaSinkers = append(notificaSinkers, &notification.WebhookSinker{Url: url})
			case config.RulerSinkTypeAlert:
				alertSinkers = append(alertSinkers, &alert.WebhookSinker{Url: url})
			}
		}
		if a := c.Sinks.Alertmanager; a != nil {
			amsinker, e := alert.NewAlertmanagerSinker(a)
			if e != nil {
				return e
			}
			alertSinkers = append(alertSinkers, amsinker)
		}
		if so := c.Sinks.Stdout; so != nil {
			switch so.Type {
			case config.RulerSinkTypeNotification:
				notificaSinkers = append(notificaSinkers, &notification.StdoutSinker{})
			case config.RulerSinkTypeAlert:
				alertSinkers = append(alertSinkers, &alert.StdoutSinker{})
			}
		}
	}

	ruleCache, e := NewRuleCache(r.kcfg, c)
	if e != nil {
		return e
	}
	ctx, cancel := context.WithCancel(context.Background())
	if e = ruleCache.Run(ctx); e != nil {
		return e
	}

	if r.rulerConds != nil {
		if c := r.rulerConds.cancel; c != nil {
			c()
		}
	}

	r.rulerConds = &rulerConds{
		notificationSinkers: notificaSinkers,
		alertSinkers:        alertSinkers,
		ruleCache:           ruleCache,
		cancel:              cancel,
	}

	return nil
}

func (r *KubeEventsRuler) getRulerConds() *rulerConds {
	r.rulerCondsMutex.Lock()
	defer r.rulerCondsMutex.Unlock()
	return r.rulerConds
}

func (r *KubeEventsRuler) AddEvents(evts []*types.Event) {
	for _, evt := range evts {
		r.evtQueue.Add(evt)
	}
}

func (r *KubeEventsRuler) drainEvents() (evts []*types.Event, shutdown bool) {
	l := r.evtQueue.Len()
	if l <= 0 {
		shutdown = r.evtQueue.ShuttingDown()
		return
	}
	for i := 0; i < l; i++ {
		obj, sd := r.evtQueue.Get()
		if sd {
			shutdown = sd
			break
		}
		evts = append(evts, obj.(*types.Event))
	}
	return
}

func (r *KubeEventsRuler) drainNotifications() ([]*types.EventNotification, bool) {
	l := r.notificaQueue.Len()
	if l <= 0 {
		return nil, r.notificaQueue.ShuttingDown()
	}
	var evtNotifications []*types.EventNotification
	for i := 0; i < l; i++ {
		obj, sd := r.notificaQueue.Get()
		if sd {
			return evtNotifications, sd
		}
		evtNotifications = append(evtNotifications, obj.(*types.EventNotification))
	}
	return evtNotifications, false
}

func (r *KubeEventsRuler) drainAlerts() ([]*types.EventAlert, bool) {
	l := r.alertQueue.Len()
	if l <= 0 {
		return nil, r.alertQueue.ShuttingDown()
	}
	var evtAlerts []*types.EventAlert
	for i := 0; i < l; i++ {
		obj, sd := r.alertQueue.Get()
		if sd {
			return evtAlerts, sd
		}
		evtAlerts = append(evtAlerts, obj.(*types.EventAlert))
	}
	return evtAlerts, false
}

func (r *KubeEventsRuler) evalEvents(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		evts, shutdown := r.drainEvents()
		if shutdown {
			return
		}
		if len(evts) == 0 {
			continue
		}
		rc := r.getRulerConds()
		evalNotifica := len(rc.notificationSinkers) > 0
		evalAlert := len(rc.alertSinkers) > 0
		for _, evt := range evts {
			if err := r.taskPool.Submit(func() {
				defer r.evtQueue.Done(evt)
				rules := rc.ruleCache.GetRules(ctx, evt)
				if evalNotifica {
					n, err := evt.EvalToNotification(rules)
					if err != nil {
						klog.Error("error evaluating event", err)
						r.evtQueue.AddRateLimited(evt)
						return
					}
					if n != nil {
						r.notificaQueue.Add(n)
					}
				}
				if evalAlert {
					a, err := evt.EvalToAlert(rules)
					if err != nil {
						klog.Error("error evaluating event", err)
						r.evtQueue.AddRateLimited(evt)
						return
					}
					if a != nil {
						r.alertQueue.Add(a)
					}
				}
				r.evtQueue.Forget(evt)
			}); err != nil {
				klog.Error("error submitting task", err)
				r.evtQueue.AddRateLimited(evt)
				r.evtQueue.Done(evt)
			}
		}
	}
}

func (r *KubeEventsRuler) sinkNotifications(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		notificas, shutdown := r.drainNotifications()
		if shutdown {
			return
		}
		if len(notificas) == 0 {
			continue
		}

		func() {
			postFunc := r.notificaQueue.Forget
			defer func() {
				for _, n := range notificas {
					postFunc(n)
					r.notificaQueue.Done(n)
				}
			}()
			notificaSinkers := r.getRulerConds().notificationSinkers
			if len(notificaSinkers) == 0 {
				return
			}
			for _, sinker := range notificaSinkers {
				if e := sinker.SinkNotifications(ctx, notificas); e != nil {
					klog.Error("Error sinking notifications: ", e)
					postFunc = r.notificaQueue.AddRateLimited
					break
				}
			}
		}()
	}
}

func (r *KubeEventsRuler) sinkAlerts(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		alerts, shutdown := r.drainAlerts()
		if shutdown {
			return
		}
		if len(alerts) == 0 {
			continue
		}
		func() {
			postFunc := r.alertQueue.Forget
			defer func() {
				for _, a := range alerts {
					postFunc(a)
					r.alertQueue.Done(a)
				}
			}()
			alertSinkers := r.getRulerConds().alertSinkers
			if len(alertSinkers) == 0 {
				return
			}
			for _, sinker := range alertSinkers {
				if e := sinker.SinkAlerts(ctx, alerts); e != nil {
					klog.Error("Error sinking alerts: ", e)
					postFunc = r.alertQueue.AddRateLimited
					return
				}
			}
		}()
	}
}

func (r *KubeEventsRuler) Run(ctx context.Context) error {
	defer func() {
		r.evtQueue.ShutDown()
		r.notificaQueue.ShutDown()
		r.alertQueue.ShutDown()
	}()

	go r.evalEvents(ctx)
	go r.sinkNotifications(ctx)
	go r.sinkAlerts(ctx)

	<-ctx.Done()

	if rc := r.getRulerConds(); rc != nil {
		if cancel := rc.cancel; cancel != nil {
			cancel()
		}
	}

	return ctx.Err()
}

func NewKubeEventsRuler(cfg *rest.Config, taskPool *ants.Pool) *KubeEventsRuler {
	return &KubeEventsRuler{
		kcfg:          cfg,
		evtQueue:      workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "events"),
		notificaQueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "eventNotifications"),
		alertQueue:    workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "eventAlerts"),
		taskPool:      taskPool,
	}
}

type rulerConds struct {
	notificationSinkers []types.NotificationSinker
	alertSinkers        []types.AlertSinker
	ruleCache           *RuleCache
	cancel              context.CancelFunc
}
