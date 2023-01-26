package grafana_test

import (
	"context"
	"flag"
	"github.com/clambin/grafana-exporter/grafana"
	"github.com/clambin/grafana-exporter/grafana/mock"
	"github.com/gosimple/slug"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

var update = flag.Bool("update", false, "update .golden files")

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

	gp := filepath.Join("testdata", slug.Make(t.Name())+".golden")
	if *update {
		err = os.WriteFile(gp, []byte(content), 0644)
		require.NoError(t, err)
	}

	expected, err := os.ReadFile(gp)
	require.NoError(t, err)

	assert.True(t, ok)
	assert.Equal(t, string(expected), content)
}
