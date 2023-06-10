package commands

import (
	"github.com/clambin/grafana-exporter/internal/configmap"
	"github.com/clambin/grafana-exporter/internal/writer"
)

func ExportDashboardProvisioning(w writer.Writer, cfg Config) (err error) {
	content := writer.Directories{".": writer.Files{"dashboards.yml": []byte(dashboardProvisioning)}}
	if cfg.AsConfigMap {
		var asConfigMap []byte
		if _, asConfigMap, err = configmap.Serialize(content["."], "grafana-provisioning-dashboards", cfg.Namespace, ""); err == nil {
			content["."]["dashboards.yml"] = asConfigMap
		}
	}
	if err == nil {
		err = w.Write(content)
	}
	return err
}

const dashboardProvisioning = `apiVersion: 1
providers:
- name: 'dashboards'
  orgId: 1
  folder: ''
  disableDeletion: false
  updateIntervalSeconds: 3600
  allowUiUpdates: true
  options:
    path: /dashboards
    foldersFromFilesStructure: true
`
