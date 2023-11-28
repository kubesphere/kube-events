package types

import (
	"context"

	v1 "k8s.io/api/core/v1"
)

type Events struct {
	KubeEvents []*ExtendedEvent `json:"kubeEvents"`
}

type ExtendedEvent struct {
	*v1.Event `json:",inline"`
	Cluster   string `json:"cluster,omitempty"`
}

type Sinker interface {
	Sink(ctx context.Context, events Events) error
}
