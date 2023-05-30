package export

import (
	"context"
	"fmt"
	"github.com/clambin/grafana-exporter/configmap"
	"github.com/clambin/grafana-exporter/grafana"
	"github.com/clambin/grafana-exporter/writer"
	"golang.org/x/exp/slog"
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
		return fmt.Errorf("failed to get grafana dashboards: %w", err)
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
				slog.Error("failed to write dashboard file(s)", "err", err, "folder", folder)
				break
			}
			slog.Info("Wrote dashboard file", "folder", folder)
		}
	}

	return
}
