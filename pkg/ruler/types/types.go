package types

import (
	"context"
	"fmt"
	"strings"
	"sync"

	amkit "github.com/kubesphere/alertmanager-kit"
	visitor "github.com/kubesphere/event-rule-engine/visitor"
	eventsv1alpha1 "github.com/kubesphere/kube-events/pkg/apis/v1alpha1"
	"github.com/kubesphere/kube-events/pkg/util"
	corev1 "k8s.io/api/core/v1"
)

type Event struct {
	Event   *corev1.Event
	once    sync.Once
	flatEvt map[string]interface{}

	EnqueueAlert       bool
	EnqueueNotificaion bool
}

func (evt *Event) Flat() map[string]interface{} {
	evt.once.Do(func() {
		evt.flatEvt = util.StructToFlatMap(evt.Event, "", ".")
	})
	return evt.flatEvt
}

func (evt *Event) EvalByRule(rule *eventsv1alpha1.EventRule) (ok bool, err error) {
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("eval error: %v. event[%s/%s], rule[name:%s, condition:%s]",
				p, evt.Event.Namespace, evt.Event.Name, rule.Name, rule.Condition)
		}
	}()
	if rule.Condition != "" {
		err, ok = visitor.EventRuleEvaluate(evt.Flat(), rule.Condition)
	}
	return
}

func (evt *Event) EvalToNotification(evtRules []*eventsv1alpha1.Rule) (*EventNotification, error) {
	for _, er := range evtRules {
		for _, rule := range er.Spec.Rules {
			if rule.Enable && rule.Type == eventsv1alpha1.RuleTypeNotification {
				ok, e := evt.EvalByRule(&rule)
				if e != nil {
					return nil, e
				}
				if ok {
					return &EventNotification{Event: evt.Event}, nil
				}
			}
		}
	}
	return nil, nil
}

func (evt *Event) EvalToAlert(evtRules []*eventsv1alpha1.Rule) (*EventAlert, error) {
	for _, er := range evtRules {
		for _, rule := range er.Spec.Rules {
			if rule.Enable && rule.Type == eventsv1alpha1.RuleTypeAlert {
				ok, e := evt.EvalByRule(&rule)
				if e != nil {
					return nil, e
				}
				if ok {
					return generateAlert(evt, &rule), nil
				}
			}
		}
	}
	return nil, nil
}

func generateAlert(evt *Event, rule *eventsv1alpha1.EventRule) *EventAlert {
	alert := &amkit.RawAlert{
		Annotations: map[string]string{},
		Labels: map[string]string{
			"alertname": rule.Name,
			"alerttype": "event",
			strings.ToLower(evt.Event.InvolvedObject.Kind): evt.Event.InvolvedObject.Name,
		},
	}
	if ns := evt.Event.InvolvedObject.Namespace; ns != "" {
		alert.Labels["namespace"] = ns
	}
	fp := evt.Event.InvolvedObject.FieldPath
	if strings.HasPrefix(fp, "spec.containers") {
		alert.Labels["container"] = strings.TrimSuffix(strings.TrimPrefix(fp, "spec.containers{"), "}")
	} else if strings.HasPrefix(fp, "spec.initContainers{") {
		alert.Labels["container"] = strings.TrimSuffix(strings.TrimPrefix(fp, "spec.initContainers{"), "}")
	}
	for k, v := range rule.Labels {
		alert.Labels[k] = v
	}

	for k, v := range rule.Annotations {
		alert.Annotations[k] = util.FormatMap(v, evt.Flat())
	}

	return &EventAlert{
		Alert: alert,
	}
}

type EventNotification struct {
	Event *corev1.Event
}

type EventAlert struct {
	Alert *amkit.RawAlert
}

type NotificationSinker interface {
	SinkNotifications(ctx context.Context, evtNotifications []*EventNotification) error
}

type AlertSinker interface {
	SinkAlerts(ctx context.Context, evtAlerts []*EventAlert) error
}

type EventSource interface {
	Events() <-chan *Event
}
