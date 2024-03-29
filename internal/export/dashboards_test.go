package export_test

import (
	"fmt"
	"github.com/clambin/grafana-exporter/internal/export"
	"github.com/clambin/grafana-exporter/internal/fetcher"
	"github.com/clambin/grafana-exporter/internal/writer"
	"github.com/clambin/grafana-exporter/internal/writer/fs"
	"github.com/gosimple/slug"
	gapi "github.com/grafana/grafana-api-golang-client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

func TestExportDashboards(t *testing.T) {
	testcases := []struct {
		name      string
		cfg       export.Config
		filenames []string
	}{
		{
			name:      "configmap - bar",
			cfg:       export.Config{Namespace: "default", Folders: []string{"bar"}},
			filenames: []string{"grafana-dashboards-bar.yml"},
		},
		{
			name:      "configmap - foobar",
			cfg:       export.Config{Namespace: "default", Folders: []string{"foobar"}},
			filenames: []string{"grafana-dashboards-foobar.yml"},
		},
		{
			name:      "configmap - both",
			cfg:       export.Config{Namespace: "default", Folders: []string{"bar", "foobar"}},
			filenames: []string{"grafana-dashboards-bar.yml", "grafana-dashboards-foobar.yml"},
		},
		{
			name:      "direct - bar",
			cfg:       export.Config{Direct: true, Folders: []string{"bar"}},
			filenames: []string{"bar/foo.json"},
		},
		{
			name:      "direct - foobar",
			cfg:       export.Config{Direct: true, Folders: []string{"foobar"}},
			filenames: []string{"foobar/snafu.json"},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			tmpdir, err := os.MkdirTemp("", "")
			require.NoError(t, err)

			f := fakeDashboardClient{}
			w := writer.Writer{StorageHandler: &fs.Client{}, BaseDirectory: tmpdir}

			err = export.Dashboards(&f, &w, slog.Default(), tt.cfg)
			require.NoError(t, err)

			for _, filename := range tt.filenames {
				content, err := os.ReadFile(path.Join(tmpdir, filename))
				require.NoError(t, err)

				gp := filepath.Join("testdata", strings.ToLower(t.Name())+"-"+slug.Make(filename)+".yaml")
				if *update {
					require.NoError(t, os.WriteFile(gp, content, 0644))
				}
				golden, err := os.ReadFile(gp)
				require.NoError(t, err)
				assert.Equal(t, string(golden), string(content))
			}

			assert.NoError(t, os.RemoveAll(tmpdir))
		})
	}
}

var _ fetcher.DashboardClient = &fakeDashboardClient{}

type fakeDashboardClient struct{}

func (f fakeDashboardClient) Dashboards() ([]gapi.FolderDashboardSearchResponse, error) {
	return []gapi.FolderDashboardSearchResponse{
		{Title: "foo", Type: "dash-db", FolderTitle: "bar", UID: "1"},
		{Title: "snafu", Type: "dash-db", FolderTitle: "foobar", UID: "2"},
	}, nil
}

func (f fakeDashboardClient) DashboardByUID(uid string) (*gapi.Dashboard, error) {
	var dashboards = map[string]*gapi.Dashboard{
		"1": {Model: map[string]any{"folder": "bar", "title": "foo"}},
		"2": {Model: map[string]any{"folder": "foobar", "title": "snafu"}},
	}
	if dashboard, ok := dashboards[uid]; ok {
		return dashboard, nil
	}
	return nil, fmt.Errorf("invalid uid: %s", uid)
}
