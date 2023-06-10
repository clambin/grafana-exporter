package commands

import (
	"bytes"
	"fmt"
	"github.com/clambin/grafana-exporter/internal/configmap"
	"github.com/clambin/grafana-exporter/internal/fetcher"
	"github.com/clambin/grafana-exporter/internal/writer"
	gapi "github.com/grafana/grafana-api-golang-client"
	"gopkg.in/yaml.v3"
)

func ExportDataSources(f fetcher.DataSourcesClient, w writer.Writer, cfg Config) error {
	sources, err := f.DataSources()
	if err != nil {
		return fmt.Errorf("grafana get datasources: %w", err)
	}
	content, err := exportDataSourcesAsFiles(sources)
	if err == nil && cfg.AsConfigMap {
		var asConfigMap []byte
		if _, asConfigMap, err = configmap.Serialize(content["."], "datasources", cfg.Namespace, ""); err == nil {
			content["."]["datasources.yml"] = asConfigMap
		}
	}
	if err == nil {
		err = w.Write(content)
	}
	return err
}

type dataSources struct {
	APIVersion  int                `yaml:"apiVersion"`
	DataSources []*gapi.DataSource `yaml:"datasources"`
}

func exportDataSourcesAsFiles(sources []*gapi.DataSource) (writer.Directories, error) {
	wrapped := dataSources{
		APIVersion:  1,
		DataSources: sources,
	}
	var buf bytes.Buffer
	err := yaml.NewEncoder(&buf).Encode(wrapped)
	if err != nil {
		return nil, err
	}
	return writer.Directories{".": writer.Files{"datasources.yml": buf.Bytes()}}, nil
}
