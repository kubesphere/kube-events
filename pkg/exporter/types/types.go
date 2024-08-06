package types

import (
	"context"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Events struct {
	KubeEvents []*ExtendedEvent `json:"kubeEvents"`
}

type ExtendedEvent struct {
	Event   client.Object `json:",inline"`
	Cluster string        `json:"cluster,omitempty"`
}

type Sinker interface {
	Sink(ctx context.Context, events Events) error
}
