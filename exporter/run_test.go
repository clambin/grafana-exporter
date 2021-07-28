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

func TestRun_Dashboards(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(grafanaMock.ServerHandler))
	defer server.Close()

	w := &writerMock.Writer{}

	args := []string{
		"unittest",
		"--out", ".",
		"--url", server.URL,
		"--token", "1234",
		"--direct",
		"dashboards",
	}

	cfg, err := exporter.GetConfiguration(args, false)
	if assert.NoError(t, err) == false {
		t.Fail()
	}
	client := grafana.New(cfg.URL, cfg.Token)
	err = exporter.Run(client, w, cfg)
	if assert.NoError(t, err) == false {
		t.Fail()
	}

	_, ok := w.GetFile("folder1", "db-1-1.json")
	assert.True(t, ok)
	_, ok = w.GetFile("General", "db-0-1.json")
	assert.True(t, ok)

}

func TestRun_DashboardProvisioning(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(grafanaMock.ServerHandler))
	defer server.Close()

	w := &writerMock.Writer{}

	args := []string{
		"unittest",
		"--out", ".",
		"--url", server.URL,
		"--token", "1234",
		"--direct",
		"dashboard-provisioning",
	}

	cfg, err := exporter.GetConfiguration(args, false)
	if assert.NoError(t, err) == false {
		t.Fail()
	}
	client := grafana.New(cfg.URL, cfg.Token)
	err = exporter.Run(client, w, cfg)
	if assert.NoError(t, err) == false {
		t.Fail()
	}

	_, ok := w.GetFile(".", "dashboards.yml")
	assert.True(t, ok)
}

func TestRun_DataSources(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(grafanaMock.ServerHandler))
	defer server.Close()

	w := &writerMock.Writer{}

	args := []string{
		"unittest",
		"--out", ".",
		"--url", server.URL,
		"--token", "1234",
		"--direct",
		"datasources",
	}

	cfg, err := exporter.GetConfiguration(args, false)
	if assert.NoError(t, err) == false {
		t.Fail()
	}
	client := grafana.New(cfg.URL, cfg.Token)
	err = exporter.Run(client, w, cfg)
	if assert.NoError(t, err) == false {
		t.Fail()
	}

	_, ok := w.GetFile(".", "datasources.yml")
	assert.True(t, ok)
}
