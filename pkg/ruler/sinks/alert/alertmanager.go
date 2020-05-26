package alert

import (
	"context"

	amkit "github.com/kubesphere/alertmanager-kit"
	"github.com/kubesphere/kube-events/pkg/config"
	"github.com/kubesphere/kube-events/pkg/ruler/types"
)

type AlertmanagerSinker struct {
	cli *amkit.AlertmanagerClient
}

func (s *AlertmanagerSinker) SinkAlerts(ctx context.Context, evtAlerts []*types.EventAlert) error {
	var alerts []*amkit.RawAlert
	for _, ea := range evtAlerts {
		alerts = append(alerts, ea.Alert)
	}
	if len(alerts) == 0 {
		return nil
	}
	return s.cli.PostAlerts(ctx, alerts)
}

func NewAlertmanagerSinker(c *config.RulerAlertmanagerSink) (*AlertmanagerSinker, error) {
	s := &AlertmanagerSinker{}
	cc := amkit.ClientConfig{
		Service: &amkit.ServiceReference{
			Namespace:  c.Namespace,
			Name:       c.Name,
			Port:       c.Port,
			TargetPort: c.TargetPort,
		},
	}
	cli, e := amkit.NewClient(cc)
	if e != nil {
		return nil, e
	}
	s.cli = cli
	return s, nil
}
