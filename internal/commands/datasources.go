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

func ExportDataSources(f fetcher.DataSourcesClient, w *writer.Writer, cfg Config) error {
	sources, err := f.DataSources()
	if err != nil {
		return fmt.Errorf("grafana get datasources: %w", err)
	}

	content, err := exportDataSourcesAsFile(sources)
	if err == nil && cfg.AsConfigMap {
		_, content, err = configmap.Serialize(map[string][]byte{"datasources.yml": content}, "datasources", cfg.Namespace, "")
	}
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}
	if err = w.Initialize(); err != nil {
		return fmt.Errorf("write init: %w", err)
	}
	if err = w.AddFile("datasources.yml", content); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return w.Store()
}

type dataSources struct {
	APIVersion  int                `yaml:"apiVersion"`
	DataSources []*gapi.DataSource `yaml:"datasources"`
}

func exportDataSourcesAsFile(sources []*gapi.DataSource) ([]byte, error) {
	wrapped := dataSources{
		APIVersion:  1,
		DataSources: sources,
	}
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(wrapped); err != nil {
		return nil, fmt.Errorf("encode: %w", err)
	}
	return buf.Bytes(), nil
}
