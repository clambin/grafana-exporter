package cli

import (
	"fmt"
	"github.com/clambin/grafana-exporter/internal/export"
	gapi "github.com/grafana/grafana-api-golang-client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
)

var (
	datasourcesCmd = &cobra.Command{
		Use:   "datasources",
		Short: "export Grafana data sources provisioning",
		RunE:  ExportDataSources,
	}
)

func ExportDataSources(_ *cobra.Command, _ []string) error {
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
	return export.ExportDataSources(c, w, export.Config{
		AsConfigMap: !viper.GetBool("direct"),
		Namespace:   viper.GetString("namespace"),
	})
}
