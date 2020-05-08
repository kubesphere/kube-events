package sinks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kubesphere/kube-events/pkg/util"
	v1 "k8s.io/api/core/v1"
)

type WebhookSinker struct {
	Url string
}

func (s *WebhookSinker) Sink(ctx context.Context, evts []*v1.Event) error {
	var buf bytes.Buffer
	for _, evt := range evts {
		if bs, err := json.Marshal(evt); err != nil {
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
