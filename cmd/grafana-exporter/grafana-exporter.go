package main

import (
	"fmt"
	"github.com/clambin/grafana-exporter/internal/commands"
	"github.com/clambin/grafana-exporter/internal/writer"
	"github.com/clambin/grafana-exporter/version"
	"github.com/grafana-tools/sdk"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
	"strings"
)

var (
	configFilename string
	rootCmd        = &cobra.Command{
		Use:   "grafana-exporter",
		Short: "exports Grafana dashboards, dashboard provisioning & datasource provisioning",
	}
	dashboardsCmd = &cobra.Command{
		Use:   "dashboards",
		Short: "export Grafana dashboards",
		RunE:  ExportDashboards,
	}
	dashboardsProvisioningCmd = &cobra.Command{
		Use:   "dashboards-provisioning",
		Short: "export Grafana dashboard provisioning",
		RunE:  ExportDashboardProvisioning,
	}
	datasourcesCmd = &cobra.Command{
		Use:   "datasources",
		Short: "export Grafana data sources provisioning",
		RunE:  ExportDataSources,
	}
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func ExportDashboards(cmd *cobra.Command, _ []string) error {
	w := writer.NewDiskWriter(viper.GetString("out"))
	c, err := sdk.NewClient(viper.GetString("url"), viper.GetString("token"), http.DefaultClient)
	if err != nil {
		return fmt.Errorf("grafana connect: %w", err)
	}

	return commands.ExportDashboards(cmd.Context(), c, w, commands.Config{
		AsConfigMap: !viper.GetBool("direct"),
		Namespace:   viper.GetString("namespace"),
		Folders:     strings.Split(viper.GetString("folders"), ","),
	})
}

func ExportDashboardProvisioning(_ *cobra.Command, _ []string) error {
	w := writer.NewDiskWriter(viper.GetString("out"))
	return commands.ExportDashboardProvisioning(w, commands.Config{
		AsConfigMap: !viper.GetBool("direct"),
		Namespace:   viper.GetString("namespace"),
	})
}

func ExportDataSources(cmd *cobra.Command, _ []string) error {
	w := writer.NewDiskWriter(viper.GetString("out"))
	c, err := sdk.NewClient(viper.GetString("url"), viper.GetString("token"), http.DefaultClient)
	if err != nil {
		return fmt.Errorf("grafana connect: %w", err)
	}
	return commands.ExportDataSources(cmd.Context(), c, w, commands.Config{
		AsConfigMap: !viper.GetBool("direct"),
		Namespace:   viper.GetString("namespace"),
	})
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(dashboardsCmd)
	rootCmd.AddCommand(dashboardsProvisioningCmd)
	rootCmd.AddCommand(datasourcesCmd)

	rootCmd.Version = version.BuildVersion

	rootCmd.PersistentFlags().StringVarP(&configFilename, "config", "c", "", "Configuration file")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Log debug messages")
	_ = viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	rootCmd.PersistentFlags().Bool("direct", false, "Write config files directory (default: write as k8s config maps)")
	_ = viper.BindPFlag("direct", rootCmd.PersistentFlags().Lookup("direct"))
	rootCmd.PersistentFlags().StringP("namespace", "n", "monitoring", "Namespace for K8s config maps")
	_ = viper.BindPFlag("namespace", rootCmd.PersistentFlags().Lookup("namespace"))

	// TODO: ideally these should only be adding to dashboardsCmd & datasourcesCmd
	rootCmd.PersistentFlags().StringP("url", "u", "", "Grafana URL")
	_ = viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url"))
	rootCmd.PersistentFlags().StringP("token", "t", "", "Grafana API token (must have admin access)")
	_ = viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
	rootCmd.PersistentFlags().StringP("out", "o", "", "Output directory")
	_ = viper.BindPFlag("out", rootCmd.PersistentFlags().Lookup("out"))

	// ctx := context.WithValue(context.Background(), "logger", slog.Default())
	//rootCmd.SetContext(ctx)

	dashboardsCmd.PersistentFlags().StringP("folders", "f", "k3s", "comma-separared list of folders to export")
	_ = viper.BindPFlag("folders", rootCmd.PersistentFlags().Lookup("folders"))
}

func initConfig() {
	if configFilename != "" {
		viper.SetConfigFile(configFilename)
	} else {
		viper.AddConfigPath("/etc/grafana-exporter/")
		viper.AddConfigPath("$HOME/.grafana-exporter")
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
	}

	viper.SetEnvPrefix("GRAFANA_EXPORTER")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		slog.Error("failed to read config file", "err", err)
	}

	var opts slog.HandlerOptions
	if viper.GetBool("debug") {
		opts.Level = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &opts)))
}
