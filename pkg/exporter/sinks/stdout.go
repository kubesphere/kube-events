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

		eventBytes, err := json.Marshal(evt.Event)
		if err != nil {
			return err
		}

		var eventMap map[string]interface{}
		if err := json.Unmarshal(eventBytes, &eventMap); err != nil {
			return err
		}

		eventMap["cluster"] = evt.Cluster

		finalBytes, err := json.Marshal(eventMap)
		if err != nil {
			return fmt.Errorf("failed to marshal final map: %w", err)
		}

		fmt.Println(string(finalBytes))
	}
	return nil
}
