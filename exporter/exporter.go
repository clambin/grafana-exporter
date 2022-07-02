package exporter

import (
	"fmt"
	"github.com/clambin/grafana-exporter/export"
	"github.com/clambin/grafana-exporter/grafana"
	"github.com/clambin/grafana-exporter/writer"
)

type Exporter struct {
	GrafanaClient *grafana.Client
	Writer        writer.Writer
	Cfg           *Configuration
}

func NewFromArgs(args []string, showErrors bool) (*Exporter, error) {
	cfg, err := getConfiguration(args, showErrors)
	if err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	return &Exporter{
		GrafanaClient: grafana.New(cfg.URL, cfg.Token),
		Writer:        writer.NewDiskWriter(cfg.Out),
		Cfg:           cfg,
	}, nil
}

func (e *Exporter) Run() (err error) {
	switch e.Cfg.Command {
	case cmdDataSources:
		err = export.DataSources(e.GrafanaClient, e.Writer, e.Cfg.Direct, e.Cfg.Namespace)
	case cmdDashboardProvisioning:
		err = export.DashboardProvisioning(e.Writer, e.Cfg.Direct, e.Cfg.Namespace)
	case cmdDashboards:
		err = export.Dashboards(e.GrafanaClient, e.Writer, e.Cfg.Direct, e.Cfg.Namespace, e.Cfg.Folders)
	}

	return
}
