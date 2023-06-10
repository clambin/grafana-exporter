package commands_test

import (
	"github.com/clambin/grafana-exporter/internal/commands"
	"github.com/clambin/grafana-exporter/internal/fetcher"
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
		name      string
		cfg       commands.Config
		directory string
		filename  string
	}{
		{
			name:      "direct",
			directory: ".",
			filename:  "datasources.yml",
		},
		{
			name:      "configmap",
			cfg:       commands.Config{AsConfigMap: true, Namespace: "default"},
			directory: ".",
			filename:  "datasources.yml",
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			f := fakeDataSourcesClient{}
			w := fakeWriter{}
			err := commands.ExportDataSources(&f, &w, tt.cfg)
			require.NoError(t, err)
			require.Contains(t, w, ".")
			require.Contains(t, w["."], "datasources.yml")

			gp := path.Join("testdata", strings.ToLower(t.Name())+".yaml")
			if *update {
				require.NoError(t, os.WriteFile(gp, w["."]["datasources.yml"], 0644))
			}
			golden, err := os.ReadFile(gp)
			require.NoError(t, err)
			assert.Equal(t, string(golden), string(w["."]["datasources.yml"]))
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
