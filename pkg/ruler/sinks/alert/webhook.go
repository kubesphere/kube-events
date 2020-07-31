package alert

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/kubesphere/kube-events/pkg/util"
	"net/http"

	"github.com/kubesphere/kube-events/pkg/ruler/types"
)

type WebhookSinker struct {
	Url string
}

func (s *WebhookSinker) SinkAlerts(ctx context.Context, evtAlerts []*types.EventAlert) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(evtAlerts); err != nil {
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
	util.DrainResponse(resp)
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("error sinking to webhook(%s): bad response status: %s", s.Url, resp.Status)
	}
	return nil
}
