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

func ExportDashboards(f fetcher.DashboardClient, w writer.Writer, cfg Config) error {
	dashboards, err := fetcher.FetchDashboards(f, set.Create(cfg.Folders...))
	if err != nil {
		return fmt.Errorf("grafana get dashboards: %w", err)
	}

	files, err := exportDashboardsAsFiles(dashboards)
	if err == nil && cfg.AsConfigMap {
		files, err = wrapDashboards(files, cfg.Namespace)
	}

	if err == nil {
		err = w.Write(files)
	}

	return err
}

func exportDashboardsAsFiles(boards map[string][]fetcher.Board) (writer.Directories, error) {
	output := make(writer.Directories)
	for folder, folderContent := range boards {
		files := make(writer.Files)
		for _, file := range folderContent {
			model, err := encodeModel(file.Model)
			if err != nil {
				return nil, fmt.Errorf("encode model: %w", err)
			}
			files[slug.Make(file.Title)+".json"] = model
		}
		output[folder] = files
	}
	return output, nil
}

func encodeModel(input map[string]any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	err := enc.Encode(input)
	return buf.Bytes(), err
}

func wrapDashboards(directories writer.Directories, namespace string) (writer.Directories, error) {
	result := make(writer.Directories)
	for directory, dashboards := range directories {
		filename, content, err := configmap.Serialize(dashboards, "grafana-dashboards-"+directory, namespace, directory)
		if err != nil {
			return nil, err
		}
		current, ok := result["."]
		if !ok {
			current = make(writer.Files)
		}
		current[filename] = content
		result["."] = current
	}
	return result, nil
}
