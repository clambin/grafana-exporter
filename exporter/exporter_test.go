package exporter_test

import (
	"fmt"
	"github.com/clambin/grafana-exporter/exporter"
	grafanaMock "github.com/clambin/grafana-exporter/grafana/mock"
	writerMock "github.com/clambin/grafana-exporter/writer/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewFromArgs(t *testing.T) {
	testCases := []struct {
		args []string
		pass bool
	}{
		{
			args: []string{"unittest", "--out", "outdir", "--url", "http://localhost:8888", "--token", "GRAFANA_API_KEY", "dashboards", "--folders", "A,B,C"},
			pass: true,
		},
		{
			args: []string{"unittest", "--out", "outdir", "datasources"},
			pass: false,
		},
		{
			args: []string{"unittest", "--out", "outdir", "--url", "http://localhost:8888", "--token", "GRAFANA_API_KEY", "foo"},
			pass: false,
		},
	}

	for index, tt := range testCases {
		name := fmt.Sprintf("testcase %d", index+1)

		e, err := exporter.NewFromArgs(tt.args, false)
		if !tt.pass {
			assert.Error(t, err, name)
			continue
		}
		require.NoError(t, err, name)
		assert.NotNil(t, e, name)
	}
}

func TestRun_Dashboards(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(grafanaMock.ServerHandler))
	defer server.Close()

	args := []string{
		"unittest",
		"--out", ".",
		"--url", server.URL,
		"--token", "1234",
		"--direct",
		"dashboards",
	}

	e, err := exporter.NewFromArgs(args, false)
	require.NoError(t, err)

	w := &writerMock.Writer{}
	e.Writer = w

	err = e.Run()
	require.NoError(t, err)

	_, ok := w.GetFile("folder1", "db-1-1.json")
	assert.True(t, ok)
	_, ok = w.GetFile("General", "db-0-1.json")
	assert.True(t, ok)

}

func TestRun_DashboardProvisioning(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(grafanaMock.ServerHandler))
	defer server.Close()

	args := []string{
		"unittest",
		"--out", ".",
		"--url", server.URL,
		"--token", "1234",
		"--direct",
		"dashboard-provisioning",
	}

	e, err := exporter.NewFromArgs(args, false)
	require.NoError(t, err)

	w := &writerMock.Writer{}
	e.Writer = w

	err = e.Run()
	require.NoError(t, err)
	_, ok := w.GetFile(".", "dashboards.yml")
	assert.True(t, ok)
}

func TestRun_DataSources(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(grafanaMock.ServerHandler))
	defer server.Close()

	args := []string{
		"unittest",
		"--out", ".",
		"--url", server.URL,
		"--token", "1234",
		"--direct",
		"datasources",
	}

	e, err := exporter.NewFromArgs(args, false)
	require.NoError(t, err)

	w := &writerMock.Writer{}
	e.Writer = w

	err = e.Run()
	require.NoError(t, err)

	_, ok := w.GetFile(".", "datasources.yml")
	assert.True(t, ok)
}
