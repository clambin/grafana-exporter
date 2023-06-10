package commands_test

import (
	"context"
	"github.com/clambin/grafana-exporter/internal/commands"
	"github.com/clambin/grafana-exporter/internal/fetcher"
	"github.com/grafana-tools/sdk"
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
			err := commands.ExportDataSources(context.Background(), &f, &w, tt.cfg)
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

func (f fakeDataSourcesClient) GetAllDatasources(_ context.Context) ([]sdk.Datasource, error) {
	return []sdk.Datasource{
		{
			Name: "prometheus",
			Type: "prometheus",
			URL:  "http://monitoring-prometheus-mock.monitoring.svc:80",
		},
		{
			Name: "postgres",
			Type: "postgres",
			URL:  "http://postgres.default:5432",
			JSONData: struct {
				PostgresVersion int    `yaml:"postgresVersion"`
				SSSLMode        string `yaml:"sslmode"`
			}{
				PostgresVersion: 1200,
				SSSLMode:        "disable",
			},
		},
	}, nil
}
