package exporter_test

import (
	"github.com/clambin/grafana-exporter/exporter"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetConfiguration(t *testing.T) {
	args := []string{
		"unittest",
		"--out", "outdir",
		"--url", "http://localhost:8888",
		"--token", "GRAFANA_API_KEY",
		"dashboards",
		"--folders", "A,B,C",
	}

	cfg, err := exporter.GetConfiguration(args, true)
	assert.NoError(t, err)
	assert.Equal(t, "outdir", cfg.Out)
	assert.Equal(t, "http://localhost:8888", cfg.URL)
	assert.Equal(t, "GRAFANA_API_KEY", cfg.Token)
	assert.False(t, cfg.Direct)
	assert.Equal(t, "monitoring", cfg.Namespace)
	if assert.Len(t, cfg.Folders, 3) {
		assert.Contains(t, cfg.Folders, "A")
		assert.Contains(t, cfg.Folders, "B")
		assert.Contains(t, cfg.Folders, "C")
	}
}

func TestGetConfiguration_NoFolders(t *testing.T) {
	args := []string{
		"unittest",
		"--out", "outdir",
		"--url", "http://localhost:8888",
		"--token", "GRAFANA_API_KEY",
		"datasources",
	}

	cfg, err := exporter.GetConfiguration(args, true)
	assert.NoError(t, err)
	assert.Equal(t, "outdir", cfg.Out)
	assert.Equal(t, "http://localhost:8888", cfg.URL)
	assert.Equal(t, "GRAFANA_API_KEY", cfg.Token)
	assert.False(t, cfg.Direct)
	assert.Equal(t, "monitoring", cfg.Namespace)
	assert.Equal(t, "datasources", cfg.Command)
	assert.Len(t, cfg.Folders, 0)
}

func TestGetConfiguration_BadArguments(t *testing.T) {
	args := []string{
		"unittest",
		"--out", "outdir",
		"datasources",
	}
	_, err := exporter.GetConfiguration(args, true)
	assert.Error(t, err)

	args = []string{
		"unittest",
		"--out", "outdir",
		"--url", "http://localhost:8888",
		"--token", "GRAFANA_API_KEY",
		"foo",
	}
	_, err = exporter.GetConfiguration(args, false)
	assert.Error(t, err)
	assert.Equal(t, "expected command but got \"foo\"", err.Error())
}
