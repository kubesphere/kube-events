package exporter

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/kubesphere/kube-events/pkg/util"
	eventsv1 "k8s.io/api/events/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kubesphere/kube-events/pkg/config"
	"github.com/kubesphere/kube-events/pkg/exporter/sinks"
	"github.com/kubesphere/kube-events/pkg/exporter/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
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

type K8sEventSource struct {
	client    *kubernetes.Clientset
	workqueue workqueue.RateLimitingInterface
	inf       cache.SharedIndexInformer
	sinkers   []types.Sinker
	mutex     sync.Mutex

	cluster string
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

func (s *K8sEventSource) getClusterName() string {
	ns, err := s.client.CoreV1().Namespaces().Get(context.Background(), "kubesphere-system", metav1.GetOptions{})
	if err != nil {
		klog.Errorf("get namespace kubesphere-system error: %s", err)
		return ""
	}

	if ns.Annotations != nil {
		return ns.Annotations["cluster.kubesphere.io/name"]
	}

	return ""
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
		return errors.New("failed to sync events cache")
	}
	klog.Info("Successfully synced events cache")
	return nil
}

func (s *K8sEventSource) drainEvents() (evts []client.Object, shutdown bool) {
	var (
		i = 0
		m = s.workqueue.Len()
	)
	if m > maxBatchSize {
		m = maxBatchSize
	}
	for {
		var obj interface{}
		obj, shutdown = s.workqueue.Get()
		if obj != nil {
			evts = append(evts, obj.(client.Object))
		}
		i++
		if i >= m {
			break
		}
	}
	return
}

func (s *K8sEventSource) sinkEvents(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		evts, shutdown := s.drainEvents()
		if len(evts) == 0 {
			if shutdown {
				return
			}
			continue
		}

		func() {
			var err error
			defer func() {
				for _, evt := range evts {
					if err == nil {
						s.workqueue.Forget(evt)
					} else if numRequeues := s.workqueue.NumRequeues(evt); numRequeues >= maxRetries {
						s.workqueue.Forget(evt)
						klog.Infof("Dropping event %s/%s out of the queue because of failing %d times: %v\n",
							evt.GetNamespace(), evt.GetName(), numRequeues, err)
					} else {
						s.workqueue.AddRateLimited(evt)
					}
					s.workqueue.Done(evt)
				}
			}()

			events := types.Events{}
			for _, e := range evts {
				events.KubeEvents = append(events.KubeEvents, &types.ExtendedEvent{
					Event:   e,
					Cluster: util.GetCluster(),
				})
			}

			evtSinkers := s.getSinkers()
			if len(evtSinkers) == 0 {
				return
			}
			for _, sinker := range evtSinkers {
				if err = sinker.Sink(ctx, events); err != nil {
					err = fmt.Errorf("error sinking events: %v", err)
					klog.Error(err)
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

	evt, ok := obj.(client.Object)
	if ok {
		evt.SetManagedFields(nil) // set it nil because it is quite verbose
		s.workqueue.Add(evt)
	}
}

func NewKubeEventSource(client *kubernetes.Clientset) *K8sEventSource {
	s := &K8sEventSource{
		client:    client,
		workqueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "events"),
	}
	var eventType runtime.Object
	var lw *cache.ListWatch
	if util.EventVersion == util.EventVersionNew {
		eventType = &eventsv1.Event{}
		lw = cache.NewListWatchFromClient(client.EventsV1().RESTClient(),
			"events", metav1.NamespaceAll, fields.Everything())
	} else {
		eventType = &corev1.Event{}
		lw = cache.NewListWatchFromClient(client.CoreV1().RESTClient(),
			"events", metav1.NamespaceAll, fields.Everything())
	}
	s.inf = cache.NewSharedIndexInformer(lw, eventType, 0, cache.Indexers{})
	s.inf.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: s.enqueueEvent,
		UpdateFunc: func(old, new interface{}) {
			s.enqueueEvent(new)
		},
	})

	s.cluster = util.GetCluster()

	return s
}
