package grafana_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"grafana_exporter/internal/grafanatest"
)

func TestGetDashboardFolders(t *testing.T) {
	var (
		dashboardMap map[string]map[string]string
		folder       map[string]string
		content      string
		ok           bool
		err          error
	)
	exportedFolders := make([]string, 0)

	dashboardMap, err = grafanatest.NewWithHTTPClient().GetAllDashboards(exportedFolders)

	if assert.Nil(t, err) {
		assert.Len(t, dashboardMap, 2)
		folder, ok = dashboardMap["General"]
		assert.True(t, ok)
		assert.Len(t, folder, 1)
		content, ok = folder["db-0-1.json"]
		assert.True(t, ok)
		assert.Equal(t, `"dashboard 2"`, content)
		folder, ok = dashboardMap["folder1"]
		assert.True(t, ok)
		assert.Len(t, folder, 1)
		content, ok = folder["db-1-1.json"]
		assert.True(t, ok)
		assert.Equal(t, `"dashboard 1"`, content)
	}

	exportedFolders = []string{"folder1"}

	dashboardMap, err = grafanatest.NewWithHTTPClient().GetAllDashboards(exportedFolders)

	if assert.Nil(t, err) {
		assert.Len(t, dashboardMap, 1)
		folder, ok = dashboardMap["General"]
		assert.False(t, ok)
		folder, ok = dashboardMap["folder1"]
		assert.True(t, ok)
		assert.Len(t, folder, 1)
		content, ok = folder["db-1-1.json"]
		assert.True(t, ok)
		assert.Equal(t, `"dashboard 1"`, content)
	}

}

func TestGetDatasources(t *testing.T) {
	datasourceMap, err := grafanatest.NewWithHTTPClient().GetDatasources()

	var ok bool
	assert.Nil(t, err)
	assert.Len(t, datasourceMap, 1)
	content, ok := datasourceMap["datasources.yml"]
	assert.True(t, ok)
	assert.Equal(t, expected, content)
}

const (
	expected = `apiVersion: 1
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
)
