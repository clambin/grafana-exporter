package exporter

import (
	"github.com/clambin/grafana-exporter/internal/configmap"
	"github.com/clambin/grafana-exporter/internal/grafana"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
)

// Configuration holds all configuration parameters
type Configuration struct {
	Debug     bool
	Direct    bool
	URL       string
	APIToken  string
	Directory string
	Namespace string
	Folders   []string
}

type writeFunc func(string, string, string) error

// Exporter exports all required data from Grafana to disk
type Exporter struct {
	configuration *Configuration
	client        *grafana.Client
	write         writeFunc
}

// New creates a new Exporter
func New(configuration *Configuration) *Exporter {
	return NewInternal(configuration, grafana.New(configuration.URL, configuration.APIToken), writeFile)
}

// NewInternal creates a new Exporter with provided Logger & Grafana Client
// Used in unit tests to test what was written to disk
func NewInternal(configuration *Configuration, client *grafana.Client, writeFunc writeFunc) *Exporter {
	return &Exporter{
		configuration: configuration,
		client:        client,
		write:         writeFunc,
	}
}

// Export writes all dashboard & datasource provisioning files to disk
func (exporter *Exporter) Export() (err error) {
	err = exporter.exportDatasourcesProvisioning()

	if err == nil {
		err = exporter.exportDashboardsProvisioning()
	}

	if err == nil {
		err = exporter.ExportDashboards()
	}

	return
}

// exportDatasources writes the Grafana datasource provisioning file
// as created in the grafana.module
func (exporter *Exporter) exportDatasourcesProvisioning() (err error) {
	var datasources map[string]string

	if datasources, err = exporter.client.GetDatasources(); err == nil {
		err = exporter.writeFiles(".", datasources, "grafana-provisioning-datasources")
	}
	return
}

// exportDashboardsProvisioning writes the Grafana dashboard provisioning file
func (exporter *Exporter) exportDashboardsProvisioning() error {
	const dashboardProvisioning = `apiVersion: 1
providers:
- name: 'dashboards'
  orgId: 1
  folder: ''
  disableDeletion: false
  updateIntervalSeconds: 3600
  allowUiUpdates: true
  options:
    path: /dashboards
    foldersFromFilesStructure: true
`
	return exporter.writeFiles(".", map[string]string{"dashboards.yml": dashboardProvisioning}, "grafana-provisioning-dashboards")
}

// ExportDashboards writes all Grafana dashboards.
// If we're writing K8S ConfigMaps. we create one config map per Grafana folder,
// each containing the JSON models as individual files.
//
// Inside the cluster, we mount each config map in a directory per folder. Using
// 'foldersFromFilesStructure: True' inside the dashboard provisioning file then
// respects that folder structure within Grafana
func (exporter *Exporter) ExportDashboards() (err error) {
	var folders map[string]map[string]string

	// get dashboards by folder
	if folders, err = exporter.client.GetAllDashboards(exporter.configuration.Folders); err == nil {
		for directory, files := range folders {
			// write all dashboards for that folder to files or a configmap
			err = exporter.writeFiles(directory, files, "grafana-dashboards-"+directory)
			if err != nil {
				break
			}
		}
	}
	return
}

func (exporter *Exporter) writeFiles(directory string, files map[string]string, configmapName string) (err error) {
	var fileName, fileContents string

	if exporter.configuration.Direct {
		targetDir := path.Join(exporter.configuration.Directory, directory)
		for fileName, fileContents = range files {
			if err = exporter.write(targetDir, fileName, fileContents); err == nil {
				log.Info("Wrote file " + path.Join(targetDir, fileName))
			} else {
				break
			}
		}
	} else {
		fileName, fileContents, err = configmap.Serialize(
			configmapName, exporter.configuration.Namespace, files)
		if err == nil {
			err = exporter.write(exporter.configuration.Directory, fileName, fileContents)
		}
		if err == nil {
			log.Info("Wrote file " + fileName)
		}
	}
	return
}

func writeFile(directory, filename string, content string) (err error) {
	err = os.MkdirAll(directory, 0755)
	if err == nil {
		err = ioutil.WriteFile(path.Join(directory, filename), []byte(content), 0644)
	}
	if err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"filename": filename,
		}).Error("unable to write file")
	}
	return
}
