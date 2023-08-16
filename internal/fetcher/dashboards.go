package fetcher

import (
	"fmt"
	"github.com/clambin/go-common/set"
	gapi "github.com/grafana/grafana-api-golang-client"
	"log/slog"
)

type DashboardClient interface {
	Dashboards() ([]gapi.FolderDashboardSearchResponse, error)
	DashboardByUID(uid string) (*gapi.Dashboard, error)
}

type Board struct {
	Title string
	Model map[string]any
}

func FetchDashboards(c DashboardClient, foldersToExport set.Set[string]) (map[string][]Board, error) {
	foundBoards, err := c.Dashboards()
	if err != nil {
		return nil, fmt.Errorf("grafana search: %w", err)
	}

	result := make(map[string][]Board)
	for _, board := range foundBoards {
		slog.Debug("dashboard found", "title", board.Title, "type", board.Type, "folder", board.FolderTitle)

		// Only process dashboards, not folders
		if board.Type != "dash-db" {
			slog.Debug("invalid type in dashboard. ignoring", "type", board.Type)
			continue
		}

		// Only export if the dashboard is in a specified folder
		if len(foldersToExport) > 0 && !foldersToExport.Contains(board.FolderTitle) {
			slog.Debug("folder not in scope. ignoring", "folderTitle", board.FolderTitle, "title", board.Title)
			continue
		}

		// Get the dashboard model
		rawBoard, err := c.DashboardByUID(board.UID)
		if err != nil {
			return nil, fmt.Errorf("grafana get board: %w", err)
		}

		boards := append(result[board.FolderTitle], Board{Title: board.Title, Model: rawBoard.Model})
		result[board.FolderTitle] = boards
	}

	return result, nil
}
