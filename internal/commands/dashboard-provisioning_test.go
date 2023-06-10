package commands_test

import (
	"github.com/clambin/grafana-exporter/internal/commands"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExportDashboardProvisioning(t *testing.T) {
	testcases := []struct {
		name      string
		config    commands.Config
		directory string
		filename  string
	}{
		{
			name:      "direct",
			directory: ".",
			filename:  "dashboards.yml",
		},
		{
			name:      "configmap",
			config:    commands.Config{AsConfigMap: true, Namespace: "default"},
			directory: ".",
			filename:  "dashboards.yml",
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			w := fakeWriter{}
			require.NoError(t, commands.ExportDashboardProvisioning(&w, tt.config))
			require.Contains(t, w, tt.directory)
			require.Contains(t, w[tt.directory], tt.filename)

			gp := filepath.Join("testdata", strings.ToLower(t.Name())+".yaml")
			if *update {
				require.NoError(t, os.WriteFile(gp, w[tt.directory][tt.filename], 0644))
			}
			golden, err := os.ReadFile(gp)
			require.NoError(t, err)
			assert.Equal(t, string(golden), string(w[tt.directory][tt.filename]))
		})
	}
}
