package sinks

import (
	"context"
	"encoding/json"
	"fmt"
	"k8s.io/api/core/v1"
)

type StdoutSinker struct {
}

func (s *StdoutSinker) Sink(ctx context.Context, evts []*v1.Event) error {
	for _, evt := range evts {
		bs, err := json.Marshal(evt)
		if err != nil {
			return err
		}
		fmt.Println(string(bs))
	}
	return nil
}
