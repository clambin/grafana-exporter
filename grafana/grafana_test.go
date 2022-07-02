package grafana_test

import (
	"context"
	"github.com/clambin/grafana-exporter/grafana"
	"github.com/clambin/grafana-exporter/grafana/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetDashboardFolders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(mock.ServerHandler))
	defer server.Close()
	client := grafana.New(server.URL, "")

	ctx := context.Background()
	dashboardMap, err := client.GetAllDashboards(ctx, []string{})
	require.NoError(t, err)
	require.Len(t, dashboardMap, 2)

	var content string
	folder, ok := dashboardMap["General"]
	require.True(t, ok)
	assert.Len(t, folder, 1)
	content, ok = folder["db-0-1.json"]
	require.True(t, ok)
	assert.Equal(t, `"dashboard 2"`, content)

	folder, ok = dashboardMap["folder1"]
	require.True(t, ok)
	require.Len(t, folder, 1)
	content, ok = folder["db-1-1.json"]
	require.True(t, ok)
	assert.Equal(t, `"dashboard 1"`, content)

	dashboardMap, err = client.GetAllDashboards(ctx, []string{"folder1"})
	require.Nil(t, err)
	require.Len(t, dashboardMap, 1)

	_, ok = dashboardMap["General"]
	assert.False(t, ok)
	folder, ok = dashboardMap["folder1"]
	require.True(t, ok)
	require.Len(t, folder, 1)
	content, ok = folder["db-1-1.json"]
	require.True(t, ok)
	assert.Equal(t, `"dashboard 1"`, content)

}

func TestGetDataSources(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(mock.ServerHandler))
	defer server.Close()
	client := grafana.New(server.URL, "")

	ctx := context.Background()
	datasourceMap, err := client.GetDataSources(ctx)

	var ok bool
	require.NoError(t, err)
	require.Len(t, datasourceMap, 1)
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
    url: http://datasource.default:5000
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
    url: http://monitoring-prometheus-mock.monitoring.svc:80
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
