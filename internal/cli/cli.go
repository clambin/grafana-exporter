package cli

import (
	"github.com/clambin/go-common/charmer"
	"github.com/clambin/grafana-exporter/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log/slog"
)

var (
	configFilename string
	RootCmd        = cobra.Command{
		Use:   "grafana-exporter",
		Short: "exports Grafana dashboards & datasources",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			charmer.SetTextLogger(cmd, viper.GetBool("debug"))
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	initArgs()
}

var args = charmer.Arguments{
	"debug":         charmer.Argument{Default: false, Help: "Log debug messages"},
	"direct":        charmer.Argument{Default: false, Help: "Export Grafana dashboards as JSON (default: write as k8s config maps)"},
	"namespace":     charmer.Argument{Default: "default", Help: "Namespace for k8s config maps"},
	"grafana.url":   charmer.Argument{Default: "http://localhost:3000", Help: "Grafana URL"},
	"grafana.token": charmer.Argument{Default: "", Help: "Grafana API token (must have admin rights)"},
	"out":           charmer.Argument{Default: ".", Help: "Output directory"},
	"mode":          charmer.Argument{Default: "local", Help: "Output mode (local/git)"},
}

func initArgs() {
	RootCmd.Version = version.BuildVersion

	RootCmd.PersistentFlags().StringVarP(&configFilename, "config", "c", "", "Configuration file")
	if err := charmer.SetPersistentFlags(&RootCmd, viper.GetViper(), args); err != nil {
		panic("failed to set flags: " + err.Error())
	}

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
		slog.Warn("failed to read config file", "err", err)
	}
}
