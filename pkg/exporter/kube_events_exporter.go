package exporter

import (
	"context"
	"fmt"
	"github.com/kubesphere/kube-events/pkg/exporter/types"
	"sync"
	"time"

	"github.com/kubesphere/kube-events/pkg/config"
	"github.com/kubesphere/kube-events/pkg/exporter/sinks"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
)

type K8sEventSource struct {
	workqueue workqueue.RateLimitingInterface
	inf       cache.SharedIndexInformer
	sinkers   []types.Sinker
	mutex     sync.Mutex
}

func (s *K8sEventSource) ReloadConfig(c *config.ExporterConfig) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var sinkers []types.Sinker
	if c == nil || c.Sinks == nil {
		s.sinkers = sinkers
		return
	}

	for _, w := range c.Sinks.Webhooks {
		if w.Url != "" {
			sinkers = append(sinkers, &sinks.WebhookSinker{Url: w.Url})
		} else if w.Service != nil {
			sinkers = append(sinkers, &sinks.WebhookSinker{Url: fmt.Sprintf("http://%s.%s.svc:%d/%s",
				w.Service.Name, w.Service.Namespace, *w.Service.Port, w.Service.Path)})
		}
	}

	if so := c.Sinks.Stdout; so != nil {
		sinkers = append(sinkers, &sinks.StdoutSinker{})
	}
	s.sinkers = sinkers
}

func (s *K8sEventSource) getSinkers() []types.Sinker {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.sinkers[:]
}

func (s *K8sEventSource) Run(ctx context.Context) error {
	defer s.workqueue.ShutDown()
	go s.sinkEvents(ctx)
	go s.inf.Run(ctx.Done())
	if err := s.waitForCacheSync(ctx.Done()); err != nil {
		return err
	}

	<-ctx.Done()
	return ctx.Err()
}

func (s *K8sEventSource) waitForCacheSync(stopc <-chan struct{}) error {
	if !cache.WaitForCacheSync(stopc, s.inf.HasSynced) {
		return fmt.Errorf("Failed to sync events cache")
	}
	klog.Info("Successfully synced events cache")
	return nil
}

func (s *K8sEventSource) drainEvents() ([]*corev1.Event, bool) {
	l := s.workqueue.Len()
	if l <= 0 {
		return nil, s.workqueue.ShuttingDown()
	}
	var evts []*corev1.Event
	for i := 0; i < l; i++ {
		obj, sd := s.workqueue.Get()
		if sd {
			return evts, sd
		}
		evts = append(evts, obj.(*corev1.Event))
	}
	return evts, false
}

func (s *K8sEventSource) sinkEvents(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		evts, shutdown := s.drainEvents()
		if shutdown {
			return
		}
		if len(evts) == 0 {
			continue
		}

		func(){
			postFunc := s.workqueue.Forget
			defer func() {
				for _, evt := range evts {
					postFunc(evt)
					s.workqueue.Done(evt)
				}
			}()
			evtSinkers := s.getSinkers()
			if len(evtSinkers) == 0 {
				return
			}
			events := &types.Events{KubeEvents: evts}
			for _, sinker := range evtSinkers {
				if e := sinker.Sink(ctx, events); e != nil {
					klog.Error("Error sinking events: ", e)
					postFunc = s.workqueue.AddRateLimited
					return
				}
			}
		}()
	}
}

func (s *K8sEventSource) enqueueEvent(obj interface{}) {
	if obj == nil {
		return
	}
	evt, ok := obj.(*corev1.Event)
	if ok {
		s.workqueue.Add(evt)
	}
}

func NewKubeEventSource(client *kubernetes.Clientset) *K8sEventSource {
	s := &K8sEventSource{
		workqueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "events"),
	}
	lw := cache.NewListWatchFromClient(client.CoreV1().RESTClient(),
		"events", metav1.NamespaceAll, fields.Everything())
	s.inf = cache.NewSharedIndexInformer(lw, &corev1.Event{}, time.Minute*30, cache.Indexers{})
	s.inf.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: s.enqueueEvent,
		UpdateFunc: func(old, new interface{}) {
			s.enqueueEvent(new)
		},
	})

	return s
}
