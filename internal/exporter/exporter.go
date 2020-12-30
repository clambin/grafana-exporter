package exporter

import (
	log "github.com/sirupsen/logrus"
	"grafana_exporter/internal/configmap"
	"grafana_exporter/internal/grafana"
	"io/ioutil"
	"path"
)

// Configuration holds all configuration parameters
type Configuration struct {
	Debug     bool
	URL       string
	APIToken  string
	Directory string
	Namespace string
	Folders   []string
}

// Exporter exports all required data from Grafana to disk
type Exporter struct {
	configuration *Configuration
	client        *grafana.Client
	write         func(string, string, []byte)
}

// New creates a new Exporter
func New(configuration *Configuration) *Exporter {
	return NewInternal(configuration, grafana.New(configuration.URL, configuration.APIToken), writeFile)
}

// NewInternal creates a new Exporter with provided Logger & Grafana Client
// Used in unit tests to test what was written to disk
func NewInternal(configuration *Configuration, client *grafana.Client, writeFunc func(string, string, []byte)) *Exporter {
	return &Exporter{
		configuration: configuration,
		client:        client,
		write:         writeFunc,
	}
}

// Export writes all dashboard & datasource provisioning files to disk
func (exporter *Exporter) Export() error {
	var err error

	err = exporter.exportDatasourcesProvisioning()

	if err == nil {
		err = exporter.exportDashboardsProvisioning()
	}

	if err == nil {
		err = exporter.ExportDashboards()
	}

	return err

}

// exportDatasources writes the Grafana datasource provisioning file
// as created in the grafana.module
func (exporter *Exporter) exportDatasourcesProvisioning() error {
	var (
		err          error
		datasources  map[string]string
		fileName     string
		fileContents string
		configMap    []byte
	)

	if datasources, err = exporter.client.GetDatasources(); err == nil {
		if true {
			for fileName, fileContents = range datasources {
				exporter.write(exporter.configuration.Directory, fileName, []byte(fileContents))
				log.Info("exported datasources provisioning file: " + fileName)

			}
		}
		if true {
			if fileName, configMap, err = configmap.Serialize(
				"grafana-provisioning-datasources", exporter.configuration.Namespace, datasources); err == nil {
				exporter.write(exporter.configuration.Directory, fileName, configMap)
				log.Info("exported config map for datasources provisioning file: " + fileName)
			}
		}
	}
	return err
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
	var (
		err              error
		fileName         string
		configMap        []byte
		provisioningFile = map[string]string{
			"dashboards.yml": dashboardProvisioning,
		}
	)

	if true {
		exporter.write(exporter.configuration.Directory, "dashboards.yml", []byte(dashboardProvisioning))
		log.Info("exported dashboard provisioning file: dashboards.yml")
	}

	if fileName, configMap, err = configmap.Serialize(
		"grafana-provisioning-dashboards", exporter.configuration.Namespace, provisioningFile); err == nil {
		exporter.write(exporter.configuration.Directory, fileName, configMap)
		log.Info("exported config map for dashboard provisioning file: grafana-provisioning-dashboards.yml")
	}

	return err
}

// ExportDashboards writes all Grafana dashboards.
// If we're writing K8S ConfigMaps. we create one config map per Grafana folder,
// each containing the JSON models as individual files.
//
// Inside the cluster, we mount each config map in a directory per folder. Using
// 'foldersFromFilesStructure: True' inside the dashboard provisioning file then
// respects that folder structure within Grafana
func (exporter *Exporter) ExportDashboards() error {
	var (
		err          error
		fileName     string
		folders      map[string]map[string]string
		fileContents string
		configMap    []byte
	)

	// get dashboards by folder
	if folders, err = exporter.client.GetAllDashboards(exporter.configuration.Folders); err == nil {
		for directory, files := range folders {
			if true {
				targetDir := path.Join(exporter.configuration.Directory, directory)
				// ensure exporter.configuration.Directory / directory exists
				for fileName, fileContents = range files {
					exporter.write(targetDir, fileName, []byte(fileContents))
					log.Info("exported dashboard file " + path.Join(directory, fileName))
				}
			}
			if true {
				if fileName, configMap, err = configmap.Serialize(
					"grafana-dashboards-"+directory, exporter.configuration.Namespace, files); err == nil {
					exporter.write(exporter.configuration.Directory, fileName, configMap)
					log.Info("exported configmap for dashboard file " + fileName)
				}
			}
			if err != nil {
				break
			}
		}
	}
	return err
}

func writeFile(directory, filename string, content []byte) {
	if err := ioutil.WriteFile(path.Join(directory, filename), content, 0644); err != nil {
		panic(err)
	}
}
