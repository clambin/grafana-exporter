package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/clambin/go-common/set"
	"github.com/clambin/grafana-exporter/internal/configmap"
	"github.com/clambin/grafana-exporter/internal/fetcher"
	"github.com/clambin/grafana-exporter/internal/writer"
	"github.com/gosimple/slug"
)

func ExportDashboards(ctx context.Context, f fetcher.DashboardClient, w writer.Writer, cfg Config) error {
	dashboards, err := fetcher.FetchDashboards(ctx, f, set.Create(cfg.Folders...))
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
			// reformat single-line json to indented multi-line layout
			pretty, err := reformatJSON(file.Content)
			if err != nil {
				return nil, fmt.Errorf("reformat json: %w", err)
			}
			files[slug.Make(file.Title)+".json"] = pretty
		}
		output[folder] = files
	}
	return output, nil
}

func reformatJSON(input []byte) ([]byte, error) {
	var buf bytes.Buffer
	var err error
	var unmarshalled any
	if err = json.Unmarshal(input, &unmarshalled); err == nil {
		enc := json.NewEncoder(&buf)
		enc.SetIndent("", "  ")
		err = enc.Encode(unmarshalled)
	}
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
