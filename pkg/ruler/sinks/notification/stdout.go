package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kubesphere/kube-events/pkg/ruler/types"
)

type StdoutSinker struct {

}

func (s *StdoutSinker) SinkNotifications(ctx context.Context, evtNotifications []*types.EventNotification) error {
	for _, notification := range evtNotifications {
		bs, err := json.Marshal(notification.Event)
		if err != nil {
			return err
		}
		fmt.Println(string(bs))
	}
	return nil
}