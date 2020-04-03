package types

import (
	"context"
	"fmt"
	"strings"
	"sync"

	amkit "github.com/kubesphere/alertmanager-kit"
	loggingv1alpha1 "github.com/kubesphere/kube-events/pkg/apis/v1alpha1"
	"github.com/kubesphere/kube-events/pkg/ruler/visitor"
	"github.com/kubesphere/kube-events/pkg/util"
	corev1 "k8s.io/api/core/v1"
)

const (
	flatEventKeyPrefix       = "event"
	flatEventKeySeparator    = "."
	alertTypeLabelName       = "alerttype"
	alertTypeLabelValueEvent = "event"
)

type Event struct {
	Event   *corev1.Event
	once    sync.Once
	flatEvt map[string]interface{}
}

func (evt *Event) Flat() map[string]interface{} {
	evt.once.Do(func() {
		evt.flatEvt = util.StructToFlatMap(evt.Event, flatEventKeyPrefix, flatEventKeySeparator)
	})
	return evt.flatEvt
}

func (evt *Event) EvalByRule(rule *loggingv1alpha1.Rule) (ok bool, err error) {
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("eval error: %v. event[%s/%s], rule[name:%s, condition:%s]",
				p, evt.Event.Namespace, evt.Event.Name, rule.Name, rule.Condition)
		}
	}()
	return visitor.EventRuleEvaluate(evt.Flat(), rule.Condition), nil
}

func (evt *Event) EvalToNotification(evtRules []*loggingv1alpha1.KubeEventsRule) (*EventNotification, error) {
	for _, er := range evtRules {
		for _, rule := range er.Spec.Rules {
			if rule.Enable && rule.Type == loggingv1alpha1.RuleTypeNotification {
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

func (evt *Event) EvalToAlert(evtRules []*loggingv1alpha1.KubeEventsRule) (*EventAlert, error) {
	for _, er := range evtRules {
		for _, rule := range er.Spec.Rules {
			if rule.Enable && rule.Type == loggingv1alpha1.RuleTypeAlert {
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

func generateAlert(evt *Event, rule *loggingv1alpha1.Rule) *EventAlert {
	alert := &amkit.RawAlert{
		Annotations: map[string]string{
			"message": util.FormatMap(rule.Message, evt.Flat()),
		},
		Labels: map[string]string{
			"alertname":        rule.Name,
			alertTypeLabelName: alertTypeLabelValueEvent,
			"namespace":        evt.Event.InvolvedObject.Namespace,
			strings.ToLower(evt.Event.InvolvedObject.Kind): evt.Event.InvolvedObject.Name,
			"severity":  rule.Priority,
			"summary":   rule.Summary,
			"summaryCn": rule.SummaryCn,
		},
	}
	fp := evt.Event.InvolvedObject.FieldPath
	if strings.HasPrefix(fp, "spec.containers") {
		alert.Labels["container"] = strings.TrimSuffix(strings.TrimPrefix(fp, "spec.containers{"), "}")
	} else if strings.HasPrefix(fp, "spec.initContainers{") {
		alert.Labels["container"] = strings.TrimSuffix(strings.TrimPrefix(fp, "spec.initContainers{"), "}")
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
