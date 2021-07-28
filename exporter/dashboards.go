package exporter

import (
	"github.com/clambin/grafana-exporter/configmap"
	"github.com/clambin/grafana-exporter/grafana"
	"github.com/clambin/grafana-exporter/writer"
	log "github.com/sirupsen/logrus"
)

func Dashboards(grafanaClient *grafana.Client, writer writer.Writer, direct bool, namespace string, folders []string) (err error) {
	var allDashboards map[string]map[string]string
	allDashboards, err = grafanaClient.GetAllDashboards(folders)

	if err == nil {
		for folder, dashboards := range allDashboards {
			if direct {
				err = writer.WriteFiles(folder, dashboards)
			} else {
				var fileName, fileContents string

				fileName, fileContents, err = configmap.Serialize("grafana-dashboards-"+folder, namespace, dashboards)

				if err == nil {
					err = writer.WriteFile(".", fileName, fileContents)
				}
			}

			if err == nil {
				log.Infof("Wrote dashboard file %s", folder)
			} else {
				log.WithError(err).Errorf("failed to write dashboard file(s) for %s: %v", folder, err.Error())
				break
			}
		}
	}

	return
}
