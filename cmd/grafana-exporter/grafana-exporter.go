package main

import (
	"github.com/clambin/grafana-exporter/internal/cli"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
	"os"
)

func main() {
	var opts slog.HandlerOptions
	if viper.GetBool("debug") {
		opts.Level = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &opts)))
	slog.Debug("debug mode")

	if err := cli.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
