package commands_test

import (
	"context"
	"github.com/clambin/grafana-exporter/internal/commands"
	"github.com/gosimple/slug"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExportDashboards(t *testing.T) {
	testcases := []struct {
		name      string
		cfg       commands.Config
		directory string
		filenames []string
		content   []byte
	}{
		{
			name:      "configmap - bar",
			cfg:       commands.Config{AsConfigMap: true, Namespace: "default", Folders: []string{"bar"}},
			directory: ".",
			filenames: []string{"grafana-dashboards-bar.yml"},
		},
		{
			name:      "configmap - foobar",
			cfg:       commands.Config{AsConfigMap: true, Namespace: "default", Folders: []string{"foobar"}},
			directory: ".",
			filenames: []string{"grafana-dashboards-foobar.yml"},
		},
		{
			name:      "configmap - both",
			cfg:       commands.Config{AsConfigMap: true, Namespace: "default", Folders: []string{"bar", "foobar"}},
			directory: ".",
			filenames: []string{"grafana-dashboards-bar.yml", "grafana-dashboards-foobar.yml"},
		},
		{
			name:      "direct - bar",
			cfg:       commands.Config{Folders: []string{"bar"}},
			directory: "bar",
			filenames: []string{"foo.json"},
		},
		{
			name:      "direct - foobar",
			cfg:       commands.Config{Folders: []string{"foobar"}},
			directory: "foobar",
			filenames: []string{"snafu.json"},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			f := fakeDashboardClient{}
			w := make(fakeWriter)
			err := commands.ExportDashboards(context.Background(), &f, &w, tt.cfg)
			require.NoError(t, err)
			require.Contains(t, w, tt.directory)
			for _, filename := range tt.filenames {
				require.Contains(t, w[tt.directory], filename)

				gp := filepath.Join("testdata", strings.ToLower(t.Name())+"-"+slug.Make(filename)+".yaml")
				if *update {
					require.NoError(t, os.WriteFile(gp, w[tt.directory][filename], 0644))
				}
				golden, err := os.ReadFile(gp)
				require.NoError(t, err)
				assert.Equal(t, string(golden), string(w[tt.directory][filename]))
			}
		})
	}
}
