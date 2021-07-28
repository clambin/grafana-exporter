package export_test

import (
	"github.com/clambin/grafana-exporter/export"
	"github.com/clambin/grafana-exporter/grafana"
	grafanaMock "github.com/clambin/grafana-exporter/grafana/mock"
	writerMock "github.com/clambin/grafana-exporter/writer/mock"
	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, err)

	content, ok := writer.GetFile(".", "grafana-dashboards-general.yml")
	if assert.True(t, ok) {
		assert.Contains(t, content, "namespace: monitoring")
	}
}

func TestDashBoards_Direct(t *testing.T) {
	writer := &writerMock.Writer{}
	server := httptest.NewServer(http.HandlerFunc(grafanaMock.ServerHandler))
	defer server.Close()
	client := grafana.New(server.URL, "")

	err := export.Dashboards(client, writer, true, "monitoring", []string{})
	assert.NoError(t, err)

	content, ok := writer.GetFile("folder1", "db-1-1.json")
	assert.True(t, ok)
	assert.Contains(t, content, "dashboard 1")

	content, ok = writer.GetFile("General", "db-0-1.json")
	assert.True(t, ok)
	assert.Contains(t, content, "dashboard 2")
}

func TestDashBoards_Filtered(t *testing.T) {
	writer := &writerMock.Writer{}
	server := httptest.NewServer(http.HandlerFunc(grafanaMock.ServerHandler))
	defer server.Close()

	client := grafana.New(server.URL, "")

	err := export.Dashboards(client, writer, true, "monitoring", []string{"General"})
	assert.NoError(t, err)

	content, ok := writer.GetFile("folder1", "db-1-1.json")
	assert.False(t, ok)

	content, ok = writer.GetFile("General", "db-0-1.json")
	assert.True(t, ok)
	assert.Contains(t, content, "dashboard 2")
}
