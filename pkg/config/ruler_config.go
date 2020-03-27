package config

import (
	"io/ioutil"
	"k8s.io/apimachinery/pkg/labels"
)

import (
	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

type RulerConfig struct {
	RuleSelector labels.Selector `json:"ruleSelector,omitempty"`
	RuleNamespaceSelector labels.Selector `json:"ruleNamespaceSelector,omitempty"`
	Sinks *RulerSinks `json:"sinks,omitempty"`
}

type RulerSinks struct {
	Alertmanager *RulerAlertmanagerSink `json:"alertmanager,omitempty"`
	Webhooks []*RulerWebhookSink `json:"webhooks,omitempty"`
	Stdout *RulerStdoutSink `json:"stdout,omitempty"`
}

type RulerAlertmanagerSink struct {
	Namespace string `json:"namespace"`
	Name string `json:"name"`
	Port *int `json:"port"`
	// TargetPort is the port to access on the backend instances targeted by the service.
	// If this is not specified, the value of the 'port' field is used.
	TargetPort *int `json:"targetPort,omitempty"`
}

type RulerWebhookSink struct {
	Type RulerSinkType `json:"type"`
	Url string `json:"namespace"`
	Service *ServiceReference `json:"service"`
}

type RulerStdoutSink struct {
	Type RulerSinkType `json:"type"`
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