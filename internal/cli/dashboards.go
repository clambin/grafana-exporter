package cli

import (
	"fmt"
	"github.com/clambin/go-common/charmer"
	"github.com/clambin/grafana-exporter/internal/export"
	gapi "github.com/grafana/grafana-api-golang-client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"strings"
)

var (
	dashboardsCmd = &cobra.Command{
		Use:   "dashboards",
		Short: "export Grafana dashboards",
		RunE:  exportDashboards,
	}
)

func exportDashboards(cmd *cobra.Command, _ []string) error {
	l := charmer.GetLogger(cmd)
	l.Info("exporting dashboards")

	w, err := makeWriter()
	if err != nil {
		return err
	}

	c, err := gapi.New(viper.GetString("grafana.url"), gapi.Config{
		APIKey: viper.GetString("grafana.token"),
		Client: http.DefaultClient,
	})
	if err != nil {
		return fmt.Errorf("grafana connect: %w", err)
	}

	return export.Dashboards(c, w, l, export.Config{
		Direct:    viper.GetBool("direct"),
		Namespace: viper.GetString("namespace"),
		Folders:   strings.Split(viper.GetString("folders"), ","),
	})
}
