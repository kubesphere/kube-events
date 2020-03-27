package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/kubesphere/kube-events/pkg/ruler/types"
	"net/http"
)

type WebhookSinker struct {
	Url string
}

func (s *WebhookSinker) SinkNotifications(ctx context.Context, evtNotifications []*types.EventNotification) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(evtNotifications); err != nil {
		return err
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