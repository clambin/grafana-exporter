package main

import (
	"github.com/clambin/grafana-exporter/export"
	"github.com/clambin/grafana-exporter/grafana"
	"github.com/clambin/grafana-exporter/version"
	"github.com/clambin/grafana-exporter/writer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
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
		Run:   ExportDashboards,
	}
	dashboardsProvisioningCmd = &cobra.Command{
		Use:   "dashboards-provisioning",
		Short: "export Grafana dashboard provisioning",
		Run:   ExportDashboardProvisioning,
	}
	datasourcesCmd = &cobra.Command{
		Use:   "datasources",
		Short: "export Grafana data sources provisioning",
		Run:   ExportDataSources,
	}
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func ExportDashboards(cmd *cobra.Command, _ []string) {
	slog.Info("exporting Grafana dashboards", "version", cmd.Root().Version)

	err := export.Dashboards(
		grafana.New(viper.GetString("url"), viper.GetString("token")),
		writer.NewDiskWriter(viper.GetString("out")),
		viper.GetBool("direct"),
		viper.GetString("namespace"),
		strings.Split(viper.GetString("folders"), ","),
	)

	if err != nil {
		slog.Error("failed to export Grafana dashboards", err)
	} else {
		slog.Info("exported Grafana dashboards")
	}
}

func ExportDashboardProvisioning(cmd *cobra.Command, _ []string) {
	slog.Info("exporting Grafana dashboard provisioning information", "version", cmd.Root().Version)

	err := export.DashboardProvisioning(
		writer.NewDiskWriter(viper.GetString("out")),
		viper.GetBool("direct"),
		viper.GetString("namespace"),
	)
	if err != nil {
		slog.Error("failed to export Grafana dashboard provisioning", err)
	} else {
		slog.Info("exported Grafana dashboard provisioning")
	}
}

func ExportDataSources(cmd *cobra.Command, _ []string) {
	slog.Info("exporting Grafana data sources provisioning information", "version", cmd.Root().Version)

	err := export.DataSources(
		grafana.New(viper.GetString("url"), viper.GetString("token")),
		writer.NewDiskWriter(viper.GetString("out")),
		viper.GetBool("direct"),
		viper.GetString("namespace"),
	)
	if err != nil {
		slog.Error("failed to export Grafana data sources provisioning", err)
	} else {
		slog.Info("exported Grafana data sources provisioning")
	}
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
		slog.Error("failed to read config file", err)
	}

	if viper.GetBool("debug") {
		slog.SetDefault(slog.New(slog.HandlerOptions{Level: slog.LevelDebug}.NewTextHandler(os.Stderr)))
	}
}
