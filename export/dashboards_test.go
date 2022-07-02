package export_test

import (
	"github.com/clambin/grafana-exporter/export"
	"github.com/clambin/grafana-exporter/grafana"
	grafanaMock "github.com/clambin/grafana-exporter/grafana/mock"
	writerMock "github.com/clambin/grafana-exporter/writer/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestDashboards_K8s(t *testing.T) {
	writer := &writerMock.Writer{}
	server := httptest.NewServer(http.HandlerFunc(grafanaMock.ServerHandler))
	defer server.Close()
	client := grafana.New(server.URL, "")
	err := export.Dashboards(client, writer, false, "monitoring", []string{})
	require.NoError(t, err)
	contents, ok := writer.GetFile(".", "grafana-dashboards-general.yml")
	require.True(t, ok)

	gp := filepath.Join("testdata", t.Name()+".golden")
	if *update {
		err = os.WriteFile(gp, []byte(contents), 0644)
		require.NoError(t, err)
	}

	var golden []byte
	golden, err = os.ReadFile(gp)
	require.NoError(t, err)
	assert.Equal(t, string(golden), contents)
}

func TestDashBoards_Direct(t *testing.T) {
	writer := &writerMock.Writer{}
	server := httptest.NewServer(http.HandlerFunc(grafanaMock.ServerHandler))
	defer server.Close()
	client := grafana.New(server.URL, "")
	err := export.Dashboards(client, writer, true, "monitoring", []string{})
	require.NoError(t, err)

	contents, ok := writer.GetFile("folder1", "db-1-1.json")
	require.True(t, ok)

	gp := filepath.Join("testdata", t.Name()+"_folder1.golden")
	if *update {
		err = os.WriteFile(gp, []byte(contents), 0644)
		require.NoError(t, err)
	}

	var golden []byte
	golden, err = os.ReadFile(gp)
	require.NoError(t, err)
	assert.Equal(t, string(golden), contents)

	contents, ok = writer.GetFile("General", "db-0-1.json")
	require.True(t, ok)

	gp = filepath.Join("testdata", t.Name()+"_General.golden")
	if *update {
		err = os.WriteFile(gp, []byte(contents), 0644)
		require.NoError(t, err)
	}

	golden, err = os.ReadFile(gp)
	require.NoError(t, err)
	assert.Equal(t, string(golden), contents)
}

func TestDashBoards_Filtered(t *testing.T) {
	writer := &writerMock.Writer{}
	server := httptest.NewServer(http.HandlerFunc(grafanaMock.ServerHandler))
	defer server.Close()
	client := grafana.New(server.URL, "")
	err := export.Dashboards(client, writer, true, "monitoring", []string{"General"})
	require.NoError(t, err)

	contents, ok := writer.GetFile("General", "db-0-1.json")
	require.True(t, ok)

	gp := filepath.Join("testdata", t.Name()+"_General.golden")
	if *update {
		err = os.WriteFile(gp, []byte(contents), 0644)
		require.NoError(t, err)
	}

	var golden []byte
	golden, err = os.ReadFile(gp)
	require.NoError(t, err)
	assert.Equal(t, string(golden), contents)

	_, ok = writer.GetFile("folder1", "db-1-1.json")
	assert.False(t, ok)
}
