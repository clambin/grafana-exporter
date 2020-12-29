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
	getArguments()

	log.Info("grafana-exporter v" + version.BuildVersion)

	exp := exporter.New(Configuration.url, Configuration.apiToken, Configuration.directory, Configuration.namespace)

	if err := exp.Export(Configuration.folders); err != nil {
		log.Warningf("failed to export: %s", err.Error())
		os.Exit(1)
	}
}

// Configuration holds all configuration parameters
var Configuration struct {
	debug     bool
	url       string
	apiToken  string
	directory string
	namespace string
	folders   []string
}

func getArguments() {
	var folders string

	a := kingpin.New(filepath.Base(os.Args[0]), "grafana provisioning exporter")
	a.Version(version.BuildVersion)
	a.HelpFlag.Short('h')
	a.VersionFlag.Short('v')
	a.Flag("debug", "Log debug messages").Short('d').BoolVar(&Configuration.debug)
	a.Flag("url", "Grafana API URL").Short('u').Required().StringVar(&Configuration.url)
	a.Flag("token", "Grafana API Token (must have admin access)").Short('t').Required().StringVar(&Configuration.apiToken)
	a.Flag("out", "Output directory").Short('o').Default(".").StringVar(&Configuration.directory)
	a.Flag("namespace", "K8s Namespace").Short('n').Default("monitoring").StringVar(&Configuration.namespace)
	a.Flag("folders", "Comma-separated list of folders to export").Short('f').Default("").StringVar(&folders)

	if _, err := a.Parse(os.Args[1:]); err != nil {
		a.Usage(os.Args[1:])
		os.Exit(1)
	}

	if folders != "" {
		Configuration.folders = strings.Split(folders, ",")
	}

	if Configuration.debug {
		log.SetLevel(log.DebugLevel)
	}
}
