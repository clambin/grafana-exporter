package main

import (
	"github.com/clambin/grafana-exporter/exporter"
	"github.com/clambin/grafana-exporter/version"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	e, err := exporter.NewFromArgs(os.Args, true)
	if err != nil {
		log.WithError(err).Fatal("export failed")
	}

	log.WithFields(log.Fields{"version": version.BuildVersion, "command": e.Cfg.Command}).Info("grafana-export")

	if e.Cfg.Debug {
		log.SetLevel(log.DebugLevel)
	}

	if err = e.Run(); err != nil {
		log.WithError(err).Fatal("export failed")
	}
}
