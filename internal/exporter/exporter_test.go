package exporter_test

import (
	"github.com/stretchr/testify/assert"
	"grafana_exporter/internal/exporter"
	"grafana_exporter/internal/grafanatest"
	"os"
	"testing"
)

func TestExporter(t *testing.T) {
	var ok bool
	var content []byte

	log := newLogger()
	dir := os.TempDir()
	err := exporter.NewExtended(
		grafanatest.NewWithHTTPClient(),
		dir,
		"monitoring",
		log.writeFile,
	).Export()

	assert.Nil(t, err)
	assert.Len(t, log.output, 1)

	for _, files := range log.output {
		content, ok = files[`grafana-provisioning-datasources.yml`]
		assert.True(t, ok)
		assert.Equal(t, datasources, string(content))
		content, ok = files[`grafana-provisioning-dashboards.yml`]
		assert.True(t, ok)
		assert.Equal(t, dashboards, string(content))
		content, ok = files[`grafana-dashboards-general.yml`]
		assert.True(t, ok)
		assert.Equal(t, general, string(content))
		content, ok = files[`grafana-dashboards-folder1.yml`]
		assert.True(t, ok)
		assert.Equal(t, folder1, string(content))
	}
}

type logger struct {
	output map[string]map[string][]byte
}

func newLogger() *logger {
	return &logger{
		output: make(map[string]map[string][]byte),
	}
}

func (log *logger) writeFile(directory, filename string, content []byte) {
	var ok bool

	if _, ok = log.output[directory]; ok == false {
		log.output[directory] = make(map[string][]byte)
	}
	log.output[directory][filename] = content
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
  name: grafana-dashboard-provisioning
  namespace: monitoring
data:
  dashboards.yml: |
    apiVersion: 1
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
