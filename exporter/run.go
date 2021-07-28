package exporter

import (
	"github.com/clambin/grafana-exporter/export"
	"github.com/clambin/grafana-exporter/grafana"
	"github.com/clambin/grafana-exporter/writer"
)

// Run the specified configuration
func Run(client *grafana.Client, w writer.Writer, cfg *Configuration) (err error) {
	switch cfg.Command {
	case cmdDataSources:
		err = export.DataSources(client, w, cfg.Direct, cfg.Namespace)
	case cmdDashboardProvisioning:
		err = export.DashboardProvisioning(w, cfg.Direct, cfg.Namespace)
	case cmdDashboards:
		err = export.Dashboards(client, w, cfg.Direct, cfg.Namespace, cfg.Folders)
	}

	return
}
