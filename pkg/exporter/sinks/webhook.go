package sinks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kubesphere/kube-events/pkg/exporter/types"
	"github.com/kubesphere/kube-events/pkg/util"
)

type WebhookSinker struct {
	Url string
}

func (s *WebhookSinker) Sink(ctx context.Context, evts *types.Events) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(evts); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", s.Url, &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("error sinking to webhook(%s): %v", s.Url, err)
	}
	util.DrainResponse(resp)
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("error sinking to webhook(%s): bad response status: %s", s.Url, resp.Status)
	}
	return nil
}
