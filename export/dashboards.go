package export

import (
	"context"
	"github.com/clambin/grafana-exporter/configmap"
	"github.com/clambin/grafana-exporter/grafana"
	"github.com/clambin/grafana-exporter/writer"
	log "github.com/sirupsen/logrus"
)

// Dashboards exports dashboards for the specified folders
//
// If direct is true, it writes the files directly, using a directory for each folder.
// Otherwise, it creates a K8S config map for the specified namespace.
func Dashboards(grafanaClient *grafana.Client, writer writer.Writer, direct bool, namespace string, folders []string) (err error) {
	ctx := context.Background()
	var allDashboards map[string]map[string]string
	allDashboards, err = grafanaClient.GetAllDashboards(ctx, folders)

	if err != nil {
		return
	}

	for folder, dashboards := range allDashboards {
		if direct {
			err = writer.WriteFiles(folder, dashboards)
		} else {
			var fileName, fileContents string

			fileName, fileContents, err = configmap.Serialize("grafana-dashboards-"+folder, namespace, folder, dashboards)

			if err == nil {
				err = writer.WriteFile(".", fileName, fileContents)
			}
			if err != nil {
				log.WithError(err).WithField("folder", folder).Error("failed to write dashboard file(s)")
				break
			}
			log.WithField("folder", folder).Infof("Wrote dashboard file %s", folder)
		}
	}

	return
}
