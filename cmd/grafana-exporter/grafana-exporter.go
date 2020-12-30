package main

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"grafana_exporter/internal/exporter"
	"grafana_exporter/internal/version"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	configuration := getArguments()

	if configuration.Debug {
		log.SetLevel(log.DebugLevel)
	}

	log.Info("grafana-exporter v" + version.BuildVersion)

	exp := exporter.New(configuration)

	if err := exp.Export(); err != nil {
		log.Warningf("failed to export: %s", err.Error())
		os.Exit(1)
	}
}

func getArguments() *exporter.Configuration {
	var (
		configuration exporter.Configuration
		folders       string
	)

	a := kingpin.New(filepath.Base(os.Args[0]), "grafana provisioning exporter")
	a.Version(version.BuildVersion)
	a.HelpFlag.Short('h')
	a.VersionFlag.Short('v')
	a.Flag("debug", "Log debug messages").Short('d').BoolVar(&configuration.Debug)
	a.Flag("configmap", "Write to K8S ConfigMaps").BoolVar(&configuration.Configmap)
	a.Flag("url", "Grafana API URL").Short('u').Required().StringVar(&configuration.URL)
	a.Flag("token", "Grafana API Token (must have admin access)").Short('t').Required().StringVar(&configuration.APIToken)
	a.Flag("out", "Output directory").Short('o').Default(".").StringVar(&configuration.Directory)
	a.Flag("namespace", "K8s Namespace").Short('n').Default("monitoring").StringVar(&configuration.Namespace)
	a.Flag("folders", "Comma-separated list of folders to export").Short('f').Default("").StringVar(&folders)

	if _, err := a.Parse(os.Args[1:]); err != nil {
		a.Usage(os.Args[1:])
		os.Exit(1)
	}

	if folders != "" {
		configuration.Folders = strings.Split(folders, ",")
	}

	return &configuration
}
