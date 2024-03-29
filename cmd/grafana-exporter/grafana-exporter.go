package main

import (
	"github.com/clambin/grafana-exporter/internal/cli"
	"github.com/clambin/grafana-exporter/pkg/charmer"
	"github.com/spf13/viper"
	"log/slog"
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
		charmer.GetLogger(cli.RootCmd).Error("failed to run", "err", err)
		os.Exit(1)
	}
}
