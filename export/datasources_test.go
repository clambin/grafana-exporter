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

func TestDataSources_Direct(t *testing.T) {
	writer := &writerMock.Writer{}
	server := httptest.NewServer(http.HandlerFunc(grafanaMock.ServerHandler))
	defer server.Close()
	client := grafana.New(server.URL, "")
	err := export.DataSources(client, writer, true, "")
	assert.NoError(t, err)
	contents, ok := writer.GetFile(".", "datasources.yml")
	assert.True(t, ok)

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

func TestDataSources_K8S(t *testing.T) {
	writer := &writerMock.Writer{}
	server := httptest.NewServer(http.HandlerFunc(grafanaMock.ServerHandler))
	defer server.Close()
	client := grafana.New(server.URL, "")
	err := export.DataSources(client, writer, false, "")
	assert.NoError(t, err)
	contents, ok := writer.GetFile(".", "grafana-provisioning-datasources.yml")
	assert.True(t, ok)

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

func TestDataSources_Failures(t *testing.T) {
	writer := &writerMock.Writer{}
	server := httptest.NewServer(http.HandlerFunc(grafanaMock.ServerHandler))
	client := grafana.New(server.URL, "")
	server.Close()

	err := export.DataSources(client, writer, false, "")
	assert.Error(t, err)
}
