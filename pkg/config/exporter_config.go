package config

import (
	"io/ioutil"

	"github.com/pkg/errors"
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
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name,omitempty"`
	Port      *int   `json:"port,omitempty"`
	Path      string `json:"path,omitempty"`
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
