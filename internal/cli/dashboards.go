package cli

import (
	"fmt"
	"github.com/clambin/grafana-exporter/internal/commands"
	gapi "github.com/grafana/grafana-api-golang-client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"strings"
)

var (
	DashboardsCmd = &cobra.Command{
		Use:   "dashboards",
		Short: "export Grafana dashboards",
		RunE:  exportDashboards,
	}
)

func exportDashboards(_ *cobra.Command, _ []string) error {
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

	return commands.ExportDashboards(c, w, commands.Config{
		AsConfigMap: !viper.GetBool("direct"),
		Namespace:   viper.GetString("namespace"),
		Folders:     strings.Split(viper.GetString("folders"), ","),
	})
}
