package exporter

import (
	"github.com/clambin/grafana-exporter/configmap"
	"github.com/clambin/grafana-exporter/writer"
)

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

func DashboardProvisioning(writer writer.Writer, direct bool, namespace string) (err error) {
	fileName := "dashboards.yml"
	contents := dashboardProvisioning

	if direct == false {
		fileName, contents, err = configmap.Serialize("grafana-provisioning-dashboards", namespace, map[string]string{fileName: contents})
	}

	if err == nil {
		err = writer.WriteFile(".", fileName, contents)
	}

	return
}
