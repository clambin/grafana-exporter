package exporter_test

import (
	"github.com/clambin/grafana-exporter/exporter"
	"github.com/clambin/grafana-exporter/grafana"
	grafanaMock "github.com/clambin/grafana-exporter/grafana/mock"
	writerMock "github.com/clambin/grafana-exporter/writer/mock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDataSources_Direct(t *testing.T) {
	writer := &writerMock.Writer{}
	server := httptest.NewServer(http.HandlerFunc(grafanaMock.ServerHandler))
	defer server.Close()
	client := grafana.New(server.URL, "")

	err := exporter.DataSources(client, writer, true, "")
	assert.NoError(t, err)

	contents, ok := writer.GetFile(".", "datasources.yml")
	assert.True(t, ok)
	assert.Contains(t, contents, "apiVersion: 1\ndatasources:\n")
}

func TestDataSources_K8S(t *testing.T) {
	writer := &writerMock.Writer{}
	server := httptest.NewServer(http.HandlerFunc(grafanaMock.ServerHandler))
	defer server.Close()
	client := grafana.New(server.URL, "")

	err := exporter.DataSources(client, writer, false, "")
	assert.NoError(t, err)

	contents, ok := writer.GetFile(".", "grafana-provisioning-datasources.yml")
	assert.True(t, ok)
	assert.Contains(t, contents, "apiVersion: 1\n    datasources:\n")
}
