package exporter_test

import (
	"github.com/clambin/grafana-exporter/internal/exporter"
	"github.com/clambin/grafana-exporter/internal/grafanatest"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExporter(t *testing.T) {
	testCases := []struct {
		direct          bool
		expectedFiles   map[string][]string
		expectedContent []struct {
			directory string
			filename  string
			content   string
		}
	}{
		{true,
			map[string][]string{
				".":       {"datasources.yml", "dashboards.yml"},
				"folder1": {"db-1-1.json"},
				"General": {"db-0-1.json"},
			},
			[]struct {
				directory string
				filename  string
				content   string
			}{
				{"General", "db-0-1.json", `"dashboard 2"`},
				{"folder1", "db-1-1.json", `"dashboard 1"`},
			},
		},
		{false,
			map[string][]string{
				".": {
					"grafana-provisioning-datasources.yml", "grafana-provisioning-dashboards.yml",
					"grafana-dashboards-general.yml", "grafana-dashboards-folder1.yml",
				},
			},
			[]struct {
				directory string
				filename  string
				content   string
			}{
				{".", "grafana-provisioning-datasources.yml", datasources},
				{".", "grafana-provisioning-dashboards.yml", dashboards},
				{".", "grafana-dashboards-general.yml", general},
				{".", "grafana-dashboards-folder1.yml", folder1},
			},
		},
	}

	for _, testCase := range testCases {

		configuration := exporter.Configuration{
			Directory: ".",
			Direct:    testCase.direct,
			Namespace: "monitoring",
		}
		log := newLogger()
		err := exporter.NewInternal(
			&configuration,
			grafanatest.NewWithHTTPClient(),
			log.writeFile,
		).Export()

		assert.Nil(t, err)

		// Check all files were created
		assert.Len(t, log.output, len(testCase.expectedFiles))
		for directory, files := range testCase.expectedFiles {
			created, ok := log.output[directory]
			if assert.True(t, ok, directory, "missing directory: %s", directory) {
				assert.Equal(t, len(files), len(created), directory+" has incorrect nr of files")
				for _, file := range files {
					_, ok := created[file]
					assert.True(t, ok, file)
				}
			}
		}

		// Check content of json files
		for _, entry := range testCase.expectedContent {
			content, ok := log.output[entry.directory][entry.filename]
			if assert.True(t, ok, entry.filename) {
				assert.Equal(t, entry.content, content)
			}
		}
	}
}

type logger struct {
	output map[string]map[string]string
}

func newLogger() *logger {
	return &logger{
		output: make(map[string]map[string]string),
	}
}

func (log *logger) writeFile(directory, filename string, content string) (err error) {
	var ok bool

	if _, ok = log.output[directory]; ok == false {
		log.output[directory] = make(map[string]string)
	}
	log.output[directory][filename] = content

	return
}

const (
	datasources = `kind: ConfigMap
apiVersion: v1
metadata:
  name: grafana-provisioning-datasources
  namespace: monitoring
data:
  datasources.yml: |
    apiVersion: 1
    datasources:
      - id: 5
        orgid: 1
        name: foo
        type: grafana-simple-json-datasource
        access: proxy
        url: http://example.com:5000
        password: ""
        user: ""
        database: ""
        basicauth: false
        basicauthuser: null
        basicauthpassword: null
        isdefault: false
        jsondata: {}
        securejsondata: null
      - id: 3
        orgid: 1
        name: bar
        type: postgres
        access: proxy
        url: http://postgres.default:5432
        password: ""
        user: grafana
        database: bar
        basicauth: false
        basicauthuser: null
        basicauthpassword: null
        isdefault: false
        jsondata:
            postgresVersion: 1200
            sslmode: disable
        securejsondata: null
      - id: 7
        orgid: 1
        name: Prometheus
        type: prometheus
        access: proxy
        url: http://monitoring-prometheus-server.monitoring.svc:80
        password: ""
        user: ""
        database: ""
        basicauth: false
        basicauthuser: null
        basicauthpassword: null
        isdefault: true
        jsondata: {}
        securejsondata: null
`

	dashboards = `kind: ConfigMap
apiVersion: v1
metadata:
  name: grafana-provisioning-dashboards
  namespace: monitoring
data:
  dashboards.yml: |
    apiVersion: 1
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
	general = `kind: ConfigMap
apiVersion: v1
metadata:
  name: grafana-dashboards-general
  namespace: monitoring
data:
  db-0-1.json: '"dashboard 2"'
`

	folder1 = `kind: ConfigMap
apiVersion: v1
metadata:
  name: grafana-dashboards-folder1
  namespace: monitoring
data:
  db-1-1.json: '"dashboard 1"'
`
)
