package types

import (
	"context"

	v1 "k8s.io/api/core/v1"
)

type Events struct {
	KubeEvents []*v1.Event `json:"kubeEvents"`
}

type Sinker interface {
	Sink(ctx context.Context, events []*v1.Event) error
}
