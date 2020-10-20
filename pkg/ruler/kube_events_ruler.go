package ruler

import (
	"context"
	"fmt"
	"sync"

	"github.com/hashicorp/go-multierror"
	"github.com/kubesphere/kube-events/pkg/config"
	"github.com/kubesphere/kube-events/pkg/ruler/sinks/alert"
	"github.com/kubesphere/kube-events/pkg/ruler/sinks/notification"
	"github.com/kubesphere/kube-events/pkg/ruler/types"
	"github.com/panjf2000/ants/v2"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
)

const (
	// maxRetries is the number of times an object will be retried before it is dropped out of the queue.
	// With the current rate-limiter in use (5ms*2^(maxRetries-1)) the following numbers represent the times
	// an object is going to be requeued:
	//
	// 5ms, 10ms, 20ms, 40ms, 80ms, 160ms, 320ms, 640ms, 1.3s, 2.6s, 5.1s, 10.2s, 20.4s, 41s, 82s
	maxRetries = 15
)

var maxBatchSize = 500

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
	var (
		i = 0
		m = r.evtQueue.Len()
	)
	if m > maxBatchSize {
		m = maxBatchSize
	}
	for {
		var obj interface{}
		obj, shutdown = r.evtQueue.Get()
		if obj != nil {
			evts = append(evts, obj.(*types.Event))
		}
		i++
		if i >= m {
			break
		}
	}
	return
}

func (r *KubeEventsRuler) drainNotifications() (evtNotifications []*types.EventNotification, shutdown bool) {
	var (
		i = 0
		m = r.notificaQueue.Len()
	)
	if m > maxBatchSize {
		m = maxBatchSize
	}
	for {
		var obj interface{}
		obj, shutdown = r.notificaQueue.Get()
		if obj != nil {
			evtNotifications = append(evtNotifications, obj.(*types.EventNotification))
		}
		i++
		if i >= m {
			break
		}
	}
	return
}

func (r *KubeEventsRuler) drainAlerts() (evtAlerts []*types.EventAlert, shutdown bool) {
	var (
		i = 0
		m = r.alertQueue.Len()
	)
	if m > maxBatchSize {
		m = maxBatchSize
	}
	for {
		var obj interface{}
		obj, shutdown = r.alertQueue.Get()
		if obj != nil {
			evtAlerts = append(evtAlerts, obj.(*types.EventAlert))
		}
		i++
		if i >= m {
			break
		}
	}
	return
}

func (r *KubeEventsRuler) evalEvents(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		evts, shutdown := r.drainEvents()
		if len(evts) == 0 {
			if shutdown {
				return
			}
			continue
		}
		rc := r.getRulerConds()
		evalNotifica := len(rc.notificationSinkers) > 0
		evalAlert := len(rc.alertSinkers) > 0
		handle := func(evt *types.Event, err error) {
			if err == nil {
				r.evtQueue.Forget(evt)
			} else if numRequeues := r.evtQueue.NumRequeues(evt); numRequeues >= maxRetries {
				r.evtQueue.Forget(evt)
				klog.Infof("Dropping event %s/%s out of the queue because of failing %d times: %v\n",
					evt.Event.Namespace, evt.Event.Name, numRequeues, err)
			} else {
				r.evtQueue.AddRateLimited(evt)
			}
			r.evtQueue.Done(evt)
		}
		for _, evt := range evts {
			if err := r.taskPool.Submit(func() {
				var err error
				defer func() {
					handle(evt, err)
				}()
				rules := rc.ruleCache.GetRules(ctx, evt)
				if evalNotifica && !evt.NotificationEvaluated {
					var n *types.EventNotification
					n, err = evt.EvalToNotification(rules)
					if err != nil {
						err = fmt.Errorf("error evaluating event: %v", err)
						klog.Error(err)
					} else {
						evt.NotificationEvaluated = true
						if n != nil {
							r.notificaQueue.Add(n)
						}
					}
				}
				if evalAlert && !evt.AlertEvaluated {
					a, e := evt.EvalToAlert(rules)
					if e != nil {
						e = fmt.Errorf("error evaluating event: %v", e)
						klog.Error(e)
						if err == nil {
							err = e
						} else {
							err = multierror.Append(err, e)
						}
					} else {
						evt.AlertEvaluated = true
						if a != nil {
							r.alertQueue.Add(a)
						}
					}
				}
			}); err != nil {
				err = fmt.Errorf("error submitting task: %v", err)
				klog.Error(err)
				handle(evt, err)
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
		if len(notificas) == 0 {
			if shutdown {
				return
			}
			continue
		}

		func() {
			var err error
			defer func() {
				for _, n := range notificas {
					if err == nil {
						r.notificaQueue.Forget(n)
					} else if numRequeues := r.notificaQueue.NumRequeues(n); numRequeues >= maxRetries {
						r.notificaQueue.Forget(n)
						klog.Infof("Dropping notification of event %s/%s out of the queue because of failing %d times: %v\n",
							n.Event.Namespace, n.Event.Name, numRequeues, err)
					} else {
						r.notificaQueue.AddRateLimited(n)
					}
					r.notificaQueue.Done(n)
				}
			}()
			notificaSinkers := r.getRulerConds().notificationSinkers
			if len(notificaSinkers) == 0 {
				return
			}
			for _, sinker := range notificaSinkers {
				if err = sinker.SinkNotifications(ctx, notificas); err != nil {
					err = fmt.Errorf("error sinking notifications: %v", err)
					klog.Error(err)
					return
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
		if len(alerts) == 0 {
			if shutdown {
				return
			}
			continue
		}

		func() {
			var err error
			defer func() {
				for _, a := range alerts {
					if err == nil {
						r.alertQueue.Forget(a)
					} else if numRequeues := r.alertQueue.NumRequeues(a); numRequeues >= maxRetries {
						r.alertQueue.Forget(a)
						klog.Infof("Dropping alert with labels %v out of the queue because of failing %d times: %v\n",
							a.Alert.Labels, numRequeues, err)
					} else {
						r.alertQueue.AddRateLimited(a)
					}
					r.alertQueue.Done(a)
				}
			}()
			alertSinkers := r.getRulerConds().alertSinkers
			if len(alertSinkers) == 0 {
				return
			}
			for _, sinker := range alertSinkers {
				if err = sinker.SinkAlerts(ctx, alerts); err != nil {
					err = fmt.Errorf("error sinking alerts: %v", err)
					klog.Error(err)
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
