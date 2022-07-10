package export

import (
	"context"
	"fmt"
	"github.com/clambin/grafana-exporter/configmap"
	"github.com/clambin/grafana-exporter/grafana"
	"github.com/clambin/grafana-exporter/writer"
)

// DataSources creates a Grafana provisioning file for all configured Grafana data sources.  Note: Grafana does not export
// passwords, so these will need to be added manually, e.g.:
//      - id: 3
//        name: PostgreSQL
//        type: postgres
//        url: postgres.default:5432
//        database: postgres_db
//        user: postgres_user
//        secureJsonData:
//          password: "postgres_password"
//
// If direct is true, it writes the file directly.  Otherwise, it creates a K8S config map for the specified namespace.
func DataSources(grafanaClient *grafana.Client, writer writer.Writer, direct bool, namespace string) (err error) {
	var dataSources map[string]string

	ctx := context.Background()
	if dataSources, err = grafanaClient.GetDataSources(ctx); err != nil {
		return fmt.Errorf("failed to get grafana data sources: %w", err)
	}

	if direct {
		return writer.WriteFiles(".", dataSources)
	}

	var fileName, contents string
	fileName, contents, err = configmap.Serialize("grafana-provisioning-datasources", namespace, "", dataSources)

	if err == nil {
		err = writer.WriteFile(".", fileName, contents)
	}

	return
}
