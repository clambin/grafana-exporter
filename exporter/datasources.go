package exporter

import (
	"github.com/clambin/grafana-exporter/configmap"
	"github.com/clambin/grafana-exporter/grafana"
	"github.com/clambin/grafana-exporter/writer"
)

func DataSources(grafanaClient *grafana.Client, writer writer.Writer, direct bool, namespace string) (err error) {
	var datasources map[string]string

	datasources, err = grafanaClient.GetDatasources()

	if err == nil {
		if direct {
			return writer.WriteFiles(".", datasources)
		}

		var fileName, contents string
		fileName, contents, err = configmap.Serialize("grafana-provisioning-datasources", namespace, datasources)

		if err == nil {
			err = writer.WriteFile(".", fileName, contents)
		}
	}

	return
}
