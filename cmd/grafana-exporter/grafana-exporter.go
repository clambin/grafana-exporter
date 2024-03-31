package main

import (
	"github.com/clambin/go-common/charmer"
	"github.com/clambin/grafana-exporter/internal/cli"
	"os"
)

func main() {
	if err := cli.RootCmd.Execute(); err != nil {
		charmer.GetLogger(&cli.RootCmd).Error("failed to run", "err", err)
		os.Exit(1)
	}
}
