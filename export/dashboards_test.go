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
	"testing"
)

func TestDashboards_K8s(t *testing.T) {
	writer := &writerMock.Writer{}
	server := httptest.NewServer(http.HandlerFunc(grafanaMock.ServerHandler))
	defer server.Close()
	client := grafana.New(server.URL, "")

	err := export.Dashboards(client, writer, false, "monitoring", []string{})
	require.NoError(t, err)

	content, ok := writer.GetFile(".", "grafana-dashboards-general.yml")
	require.True(t, ok)
	assert.Contains(t, content, "namespace: monitoring")
}

func TestDashBoards_Direct(t *testing.T) {
	writer := &writerMock.Writer{}
	server := httptest.NewServer(http.HandlerFunc(grafanaMock.ServerHandler))
	defer server.Close()
	client := grafana.New(server.URL, "")

	err := export.Dashboards(client, writer, true, "monitoring", []string{})
	require.NoError(t, err)

	content, ok := writer.GetFile("folder1", "db-1-1.json")
	require.True(t, ok)
	assert.Contains(t, content, "dashboard 1")

	content, ok = writer.GetFile("General", "db-0-1.json")
	require.True(t, ok)
	assert.Contains(t, content, "dashboard 2")
}

func TestDashBoards_Filtered(t *testing.T) {
	writer := &writerMock.Writer{}
	server := httptest.NewServer(http.HandlerFunc(grafanaMock.ServerHandler))
	defer server.Close()

	client := grafana.New(server.URL, "")

	err := export.Dashboards(client, writer, true, "monitoring", []string{"General"})
	require.NoError(t, err)

	content, ok := writer.GetFile("General", "db-0-1.json")
	require.True(t, ok)
	assert.Contains(t, content, "dashboard 2")

	_, ok = writer.GetFile("folder1", "db-1-1.json")
	assert.False(t, ok)
}
