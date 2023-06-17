package export_test

import (
	"github.com/clambin/grafana-exporter/internal/export"
	"github.com/clambin/grafana-exporter/internal/fetcher"
	"github.com/clambin/grafana-exporter/internal/writer"
	"github.com/clambin/grafana-exporter/internal/writer/fs"
	gapi "github.com/grafana/grafana-api-golang-client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"strings"
	"testing"
)

func TestExportDataSources(t *testing.T) {
	testcases := []struct {
		name     string
		cfg      export.Config
		filename string
	}{
		{
			name:     "direct",
			filename: "datasources.yml",
		},
		{
			name:     "configmap",
			cfg:      export.Config{AsConfigMap: true, Namespace: "default"},
			filename: "datasources.yml",
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			tmpdir, err := os.MkdirTemp("", "")
			require.NoError(t, err)

			f := fakeDataSourcesClient{}
			w := writer.Writer{StorageHandler: &fs.Client{}, BaseDirectory: tmpdir}

			require.NoError(t, export.ExportDataSources(&f, &w, tt.cfg))

			content, err := os.ReadFile(path.Join(tmpdir, tt.filename))
			require.NoError(t, err)

			gp := path.Join("testdata", strings.ToLower(t.Name())+".yaml")
			if *update {
				require.NoError(t, os.WriteFile(gp, content, 0644))
			}
			golden, err := os.ReadFile(gp)
			require.NoError(t, err)
			assert.Equal(t, string(golden), string(content))

			assert.NoError(t, os.RemoveAll(tmpdir))
		})
	}

}

var _ fetcher.DataSourcesClient = &fakeDataSourcesClient{}

type fakeDataSourcesClient struct{}

func (f fakeDataSourcesClient) DataSources() ([]*gapi.DataSource, error) {
	return []*gapi.DataSource{
		{
			Name: "prometheus",
			Type: "prometheus",
			URL:  "http://monitoring-prometheus-mock.monitoring.svc:80",
		},
		{
			Name: "postgres",
			Type: "postgres",
			URL:  "http://postgres.default:5432",
			JSONData: map[string]any{
				"postgresVersion": 1200,
				"sslmode":         "disable",
			},
		},
	}, nil
}
