package cli

import (
	"github.com/clambin/grafana-exporter/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

var (
	configFilename string
	RootCmd        = &cobra.Command{
		Use:   "grafana-exporter",
		Short: "exports Grafana dashboards & datasources",
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	initArgs()
}

func initArgs() {
	RootCmd.Version = version.BuildVersion

	RootCmd.Flags().StringVarP(&configFilename, "config", "c", "", "Configuration file")
	RootCmd.PersistentFlags().BoolP("debug", "d", false, "Log debug messages")
	RootCmd.PersistentFlags().Bool("direct", false, "Write config files directory (default: write as k8s config maps)")
	RootCmd.PersistentFlags().StringP("namespace", "n", "default", "Namespace for K8s config maps")
	RootCmd.PersistentFlags().StringP("grafana.url", "u", "", "Grafana URL")
	RootCmd.PersistentFlags().StringP("grafana.token", "t", "", "Grafana API token (must have admin access)")
	RootCmd.PersistentFlags().StringP("out", "o", "", "Output directory")
	RootCmd.PersistentFlags().StringP("mode", "m", "local", "Output mode (local/git)")

	_ = viper.BindPFlag("debug", RootCmd.PersistentFlags().Lookup("debug"))
	_ = viper.BindPFlag("direct", RootCmd.PersistentFlags().Lookup("direct"))
	_ = viper.BindPFlag("namespace", RootCmd.PersistentFlags().Lookup("namespace"))
	_ = viper.BindPFlag("grafana.url", RootCmd.PersistentFlags().Lookup("grafana.url"))
	_ = viper.BindPFlag("grafana.token", RootCmd.PersistentFlags().Lookup("grafana.token"))
	_ = viper.BindPFlag("out", RootCmd.PersistentFlags().Lookup("out"))
	_ = viper.BindPFlag("mode", RootCmd.PersistentFlags().Lookup("mode"))

	dashboardsCmd.Flags().StringP("folders", "f", "", "Dashboard folders to export")
	_ = viper.BindPFlag("folders", dashboardsCmd.Flags().Lookup("folders"))

	RootCmd.AddCommand(dashboardsCmd)
	RootCmd.AddCommand(datasourcesCmd)
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
}
