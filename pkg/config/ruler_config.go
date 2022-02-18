package config

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/yaml"
)

type RulerConfig struct {
	RuleSelector          labels.Selector `json:"ruleSelector,omitempty"`
	RuleNamespaceSelector labels.Selector `json:"ruleNamespaceSelector,omitempty"`
	Sinks                 *RulerSinks     `json:"sinks,omitempty"`
}

type RulerSinks struct {
	Alertmanagers []*RulerAlertmanagerSink `json:"alertmanagers,omitempty"`
	// Alertmanager will be deprecated, please use Alertmanagers instead.
	Alertmanager *RulerAlertmanagerSink `json:"alertmanager,omitempty"`
	Webhooks     []*RulerWebhookSink    `json:"webhooks,omitempty"`
	Stdout       *RulerStdoutSink       `json:"stdout,omitempty"`
}

type RulerAlertmanagerSink struct {
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name,omitempty"`
	Port      *int   `json:"port,omitempty"`
	// TargetPort is the port to access on the backend instances targeted by the service.
	// If this is not specified, the value of the 'port' field is used.
	TargetPort *int `json:"targetPort,omitempty"`
}

type RulerWebhookSink struct {
	Type    RulerSinkType     `json:"type,omitempty"`
	Url     string            `json:"namespace,omitempty"`
	Service *ServiceReference `json:"service,omitempty"`
}

type RulerStdoutSink struct {
	Type RulerSinkType `json:"type,omitempty"`
}

type RulerSinkType string

const (
	// RulerSinkTypeNotification represents event notifications sink.
	RulerSinkTypeNotification = "notification"
	// RulerSinkTypeAlert represents alert messages sink.
	RulerSinkTypeAlert = "alert"
)

func NewRulerConfig(filename string) (*RulerConfig, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg := &RulerConfig{}
	err = yaml.UnmarshalStrict(content, cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "parsing YAML file %s", filename)
	}
	return cfg, nil
}
