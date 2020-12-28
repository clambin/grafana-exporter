package exporter

import (
	"grafana_exporter/internal/configmap"
	"grafana_exporter/internal/grafana"
	"io/ioutil"
	"path"
)

// Export writes all dashboard & datasource provisioning files to disk
func Export(url, apiToken, directory, namespace string) error {
	var err error

	if err = ExportDashboards(url, apiToken, directory, namespace); err == nil {
		err = ExportDatasources(url, apiToken, directory, namespace)
	}
	return err

}

// ExportDatasources writes to disk a ConfigMap for the Grafana datasource provisioning file,
// as created in the grafana.module
func ExportDatasources(url, apiToken, directory, namespace string) error {
	var (
		err         error
		datasources = make(map[string]string)
		folderName  string
		configMap   []byte
	)

	if datasources, err = grafana.GetDatasources(url, apiToken); err == nil {
		if folderName, configMap, err = configmap.Serialize(
			"grafana-provisioning-datasources", namespace, datasources); err == nil {
			writeFile(directory, folderName+".yml", configMap)
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
func ExportDashboards(url, apiToken, directory, namespace string) error {
	var (
		err        error
		folder     string
		folderName string
		folders    map[string]map[string]string
		dashboards map[string]string
		configMap  []byte
	)

	// write provisioning file
	if _, content, err := serializeDashboardProvisioning(namespace); err == nil {
		writeFile(directory, "grafana-provisioning-dashboards.yml", content)
	}
	// get dashboards by folder
	if folders, err = grafana.GetAllDashboards(url, apiToken); err == nil {
		// write each folder in separate configmap
		for folder, dashboards = range folders {
			if folderName, configMap, err = configmap.Serialize(
				"grafana-dashboards-"+folder, namespace, dashboards); err == nil {
				writeFile(directory, folderName+".yml", configMap)
			} else {
				break
			}
		}
	}
	return err
}

func serializeDashboardProvisioning(namespace string) (string, []byte, error) {
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
		"grafana-dashboard-provisioning", namespace,
		map[string]string{"dashboards.yml": dashboardProvisioning})
}

func writeFile(directory, filename string, content []byte) {
	if err := ioutil.WriteFile(path.Join(directory, filename), content, 0644); err != nil {
		panic(err)
	}
}
