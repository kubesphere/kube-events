package alert

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kubesphere/kube-events/pkg/ruler/types"
)

type StdoutSinker struct {

}

func (s *StdoutSinker) SinkAlerts(ctx context.Context, evtAlerts []*types.EventAlert) error {
	for _, alert := range evtAlerts {
		bs, err := json.Marshal(alert.Alert)
		if err != nil {
			return err
		}
		fmt.Println(string(bs))
	}
	return nil
}