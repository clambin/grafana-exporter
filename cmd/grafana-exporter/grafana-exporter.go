package main

import (
	"github.com/clambin/grafana-exporter/exporter"
	"github.com/clambin/grafana-exporter/grafana"
	"github.com/clambin/grafana-exporter/version"
	"github.com/clambin/grafana-exporter/writer"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	cfg, err := exporter.GetConfiguration(os.Args, true)

	if err != nil {
		os.Exit(1)
	}

	log.WithFields(log.Fields{"version": version.BuildVersion, "command": cfg.Command}).Info("grafana-export")

	if cfg.Debug {
		log.SetLevel(log.DebugLevel)
	}

	log.WithFields(log.Fields{
		"command":   cfg.Command,
		"url":       cfg.URL,
		"token":     cfg.Token,
		"out":       cfg.Out,
		"direct":    cfg.Direct,
		"namespace": cfg.Namespace,
		"folders":   cfg.Folders,
	}).Debug()

	err = exporter.Run(grafana.New(cfg.URL, cfg.Token), writer.NewDiskWriter(cfg.Out), cfg)

	if err != nil {
		log.WithError(err).Error("export failed")
	}
}
