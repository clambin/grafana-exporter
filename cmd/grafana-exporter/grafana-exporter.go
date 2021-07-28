package main

import (
	"github.com/clambin/grafana-exporter/exporter"
	"github.com/clambin/grafana-exporter/grafana"
	"github.com/clambin/grafana-exporter/version"
	"github.com/clambin/grafana-exporter/writer"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	app := kingpin.New(filepath.Base(os.Args[0]), "grafana provisioning exporter")
	debug := app.Flag("debug", "Log debug messages").Short('d').Bool()
	url := app.Flag("url", "Grafana API URL").Short('u').Required().String()
	token := app.Flag("token", "Grafana API Token (must have admin access)").Short('t').Required().String()
	out := app.Flag("out", "Output directory").Short('o').Default(".").String()
	direct := app.Flag("direct", "Write config files directory (default: write as K8S ConfigMaps").Bool()
	namespace := app.Flag("namespace", "K8s Namespace").Short('n').Default("monitoring").String()

	datasources := app.Command("datasources", "export Grafana data sources")

	dashboardProvisioning := app.Command("dashboard-provisioning", "export Grafana dashboard provisioning data")

	dashboards := app.Command("dashboards", "export Grafana dashboards")
	folders := dashboards.Flag("folders", "Comma-separated list of folders to export").Short('f').String()

	command, err := app.Parse(os.Args[1:])

	if err != nil {
		app.Usage(os.Args[1:])
		os.Exit(1)
	}

	log.WithFields(log.Fields{"version": version.BuildVersion, "command": command}).Info("grafana-exporter")

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	log.WithFields(log.Fields{
		"command":   command,
		"url":       *url,
		"token":     *token,
		"out":       *out,
		"direct":    *direct,
		"namespace": *namespace,
		"folders":   *folders,
	}).Debug()

	grafanaClient := grafana.New(*url, *token)
	w := writer.NewWriter(*out)

	switch command {
	case datasources.FullCommand():
		err = exporter.DataSources(grafanaClient, w, *direct, *namespace)
	case dashboardProvisioning.FullCommand():
		err = exporter.DashboardProvisioning(w, *direct, *namespace)
	case dashboards.FullCommand():
		var folderList []string
		if *folders != "" {
			folderList = strings.Split(*folders, ",")
		}
		err = exporter.Dashboards(grafanaClient, w, *direct, *namespace, folderList)
	}

	if err != nil {
		log.WithError(err).Error("export failed")
	}
}
