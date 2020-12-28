package main

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"grafana_exporter/internal/exporter"
	"grafana_exporter/internal/version"
	"os"
	"path/filepath"
)

func main() {
	getArguments()

	log.Info("grafana-exporter v" + version.BuildVersion)

	err := exporter.New(Configuration.url, Configuration.apiToken, Configuration.directory, Configuration.namespace).Export()

	if err != nil {
		log.Warningf("failed to export: %s", err.Error())
	}
}

// Configuration holds all configuration parameters
var Configuration struct {
	debug     bool
	url       string
	apiToken  string
	directory string
	namespace string
}

func getArguments() {
	a := kingpin.New(filepath.Base(os.Args[0]), "grafana provisioning exporter")
	a.Version(version.BuildVersion)
	a.HelpFlag.Short('h')
	a.VersionFlag.Short('v')
	a.Flag("debug", "Log debug messages").Short('d').BoolVar(&Configuration.debug)
	a.Flag("url", "Grafana API URL").Short('u').Required().StringVar(&Configuration.url)
	a.Flag("token", "Grafana API Token (must have admin access)").Short('t').Required().StringVar(&Configuration.apiToken)
	a.Flag("out", "Output directory").Short('o').Default(".").StringVar(&Configuration.directory)
	a.Flag("namespace", "K8s Namespace").Short('n').Default("monitoring").StringVar(&Configuration.namespace)

	if _, err := a.Parse(os.Args[1:]); err != nil {
		a.Usage(os.Args[1:])
		os.Exit(1)
	}

	if Configuration.debug {
		log.SetLevel(log.DebugLevel)
	}
}
