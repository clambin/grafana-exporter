package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/clambin/go-common/set"
	"github.com/clambin/grafana-exporter/internal/configmap"
	"github.com/clambin/grafana-exporter/internal/fetcher"
	"github.com/clambin/grafana-exporter/internal/writer"
	"github.com/gosimple/slug"
)

func ExportDashboards(f fetcher.DashboardClient, w *writer.Writer, cfg Config) error {
	dashboards, err := fetcher.FetchDashboards(f, set.Create(cfg.Folders...))
	if err != nil {
		return fmt.Errorf("grafana get dashboards: %w", err)
	}

	var files map[string][]byte
	if cfg.AsConfigMap {
		files, err = exportDashboardsAsConfigMaps(dashboards, cfg)
	} else {
		files, err = exportDashboardsAsFiles(dashboards)
	}
	if err != nil {
		return fmt.Errorf("export: %w", err)
	}

	if err = w.Initialize(); err != nil {
		return fmt.Errorf("write init: %w", err)
	}
	for filename, content := range files {
		if err = w.AddFile(filename, content); err != nil {
			return fmt.Errorf("write %s: %w", filename, err)
		}
	}
	return w.Store()
}

func exportDashboardsAsFiles(boards map[string][]fetcher.Board) (map[string][]byte, error) {
	result := make(map[string][]byte)
	for folder, dashboards := range boards {
		for _, file := range dashboards {
			model, err := encodeModel(file.Model)
			if err != nil {
				return nil, fmt.Errorf("encode %s/%s: %w", folder, file.Title, err)
			}
			result[slug.Make(folder)+"/"+slug.Make(file.Title)+".json"] = model
		}
	}
	return result, nil
}

func encodeModel(input map[string]any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	err := enc.Encode(input)
	return buf.Bytes(), err
}

func exportDashboardsAsConfigMaps(boards map[string][]fetcher.Board, cfg Config) (map[string][]byte, error) {
	result := make(map[string][]byte)
	for folder, dashboards := range boards {
		files := make(map[string][]byte)
		for _, file := range dashboards {
			model, err := encodeModel(file.Model)
			if err != nil {
				return nil, fmt.Errorf("encode %s/%s: %w", folder, file.Title, err)
			}
			files[slug.Make(file.Title)+".json"] = model
		}
		filename, content, err := configmap.Serialize(files, "grafana-dashboards-"+slug.Make(folder), cfg.Namespace, folder)
		if err != nil {
			return nil, fmt.Errorf("configmap %s: %w", folder, err)
		}
		result[filename] = content
	}
	return result, nil
}
