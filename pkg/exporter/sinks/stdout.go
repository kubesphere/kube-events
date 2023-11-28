package sinks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kubesphere/kube-events/pkg/exporter/types"
)

type StdoutSinker struct {
}

func (s *StdoutSinker) Sink(ctx context.Context, evts types.Events) error {
	for _, evt := range evts.KubeEvents {
		bs, err := json.Marshal(evt)
		if err != nil {
			return err
		}
		fmt.Println(string(bs))
	}
	return nil
}
