package config

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"sigs.k8s.io/yaml"
)

type ExporterConfig struct {
	Sinks *ExporterSinks `json:"sinks,omitempty"`
}

type ExporterSinks struct {
	Webhooks []*ExporterSinkWebhook `json:"webhooks,omitempty"`
	Stdout   *ExporterSinkStdout    `json:"stdout,omitempty"`
}

type ExporterSinkStdout struct {
}

type ExporterSinkWebhook struct {
	Url     string            `json:"url,omitempty"`
	Service *ServiceReference `json:"service,omitempty"`
}

type ServiceReference struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Port      *int   `json:"port"`
	Path      string `json:"path"`
}

func NewExporterConfig(filename string) (*ExporterConfig, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg := &ExporterConfig{}
	err = yaml.UnmarshalStrict(content, cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "parsing YAML file %s", filename)
	}
	return cfg, nil
}
