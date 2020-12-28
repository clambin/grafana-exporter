package exporter

import (
	log "github.com/sirupsen/logrus"
	"grafana_exporter/internal/configmap"
	"grafana_exporter/internal/grafana"
	"io/ioutil"
	"path"
)

// Exporter exports all required data from Grafana to disk
type Exporter struct {
	url       string
	apiToken  string
	directory string
	namespace string
	write     func(string, string, []byte)
}

// New creates a new Exporter
func New(url, apiToken, directory, namespace string) *Exporter {
	return NewWithLogger(url, apiToken, directory, namespace, writeFile)
}

// NewWithLogger creates a new Exporter with configured Logger
// Used in unit tests to test what was written to disk
func NewWithLogger(url, apiToken, directory, namespace string, writeFunc func(string, string, []byte)) *Exporter {
	return &Exporter{
		url:       url,
		apiToken:  apiToken,
		directory: directory,
		namespace: namespace,
		write:     writeFunc,
	}
}

// Export writes all dashboard & datasource provisioning files to disk
func (exporter *Exporter) Export() error {
	var err error

	if err = exporter.ExportDashboards(); err == nil {
		err = exporter.ExportDatasources()
	}
	return err

}

// ExportDatasources writes to disk a ConfigMap for the Grafana datasource provisioning file,
// as created in the grafana.module
func (exporter *Exporter) ExportDatasources() error {
	var (
		err         error
		datasources = make(map[string]string)
		folderName  string
		configMap   []byte
	)

	if datasources, err = grafana.GetDatasources(exporter.url, exporter.apiToken); err == nil {
		if folderName, configMap, err = configmap.Serialize(
			"grafana-provisioning-datasources", exporter.namespace, datasources); err == nil {
			filename := folderName + ".yml"
			exporter.write(exporter.directory, filename, configMap)
			log.Info("exported datasource provisioning file " + filename)
		}
	}
	return err
}

// ExportDashboards writes to disk a set of ConfigMaps for all Grafana dashboards.
// We create one config map per Grafana folder, containing the JSON models as
// individual files.
//
// Inside the cluster, we mount each config map in a directory per folder. Using
// 'foldersFromFilesStructure: True' inside the dashboard provisioning file then
// respects that folder structure within Grafana
func (exporter *Exporter) ExportDashboards() error {
	var (
		err        error
		folder     string
		folderName string
		folders    map[string]map[string]string
		dashboards map[string]string
		configMap  []byte
	)

	// write provisioning file
	if _, content, err := exporter.serializeDashboardProvisioning(); err == nil {
		exporter.write(exporter.directory, "grafana-provisioning-dashboards.yml", content)
		log.Info("exported dashboard file grafana-provisioning-dashboards.yml")
	}
	// get dashboards by folder
	if folders, err = grafana.GetAllDashboards(exporter.url, exporter.apiToken); err == nil {
		// write each folder in separate configmap
		for folder, dashboards = range folders {
			if folderName, configMap, err = configmap.Serialize(
				"grafana-dashboards-"+folder, exporter.namespace, dashboards); err == nil {
				filename := folderName + ".yml"
				exporter.write(exporter.directory, filename, configMap)
				log.Info("exported dashboard file " + filename)

			} else {
				break
			}

		}
	}
	return err
}

func (exporter *Exporter) serializeDashboardProvisioning() (string, []byte, error) {
	const dashboardProvisioning = `apiVersion: 1
  providers:
  - name: 'dashboards'
    orgId: 1
    folder: ''
    disableDeletion: False
    updataIntervalSeconds': 10
    allowUiUpdates': True
    options:
    - path: '/var/lib/grafana/dashboards'
    - foldersFromFilesStructure: True
`
	return configmap.Serialize(
		"grafana-dashboard-provisioning", exporter.namespace,
		map[string]string{"dashboards.yml": dashboardProvisioning})
}

func writeFile(directory, filename string, content []byte) {
	if err := ioutil.WriteFile(path.Join(directory, filename), content, 0644); err != nil {
		panic(err)
	}
}
