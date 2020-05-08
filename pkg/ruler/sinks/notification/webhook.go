package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kubesphere/kube-events/pkg/ruler/types"
)

type WebhookSinker struct {
	Url string
}

func (s *WebhookSinker) SinkNotifications(ctx context.Context, evtNotifications []*types.EventNotification) error {
	var buf bytes.Buffer
	for _, notification := range evtNotifications {
		if bs, err := json.Marshal(notification.Event); err != nil {
			return err
		} else if _, err := buf.Write(bs); err != nil {
			return err
		} else if err := buf.WriteByte('\n'); err != nil {
			return err
		}
	}
	req, err := http.NewRequest("POST", s.Url, &buf)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sinking to webhook(%s): %v", s.Url, err)
	}
	if resp.StatusCode/100 != 2 {
		err = fmt.Errorf("error sinking to webhook(%s): bad response status: %s", s.Url, resp.Status)
	}
	return nil
}
