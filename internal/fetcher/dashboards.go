package fetcher

import (
	"context"
	"fmt"
	"github.com/clambin/go-common/set"
	"github.com/grafana-tools/sdk"
	"golang.org/x/exp/slog"
)

type DashboardClient interface {
	Search(ctx context.Context, params ...sdk.SearchParam) ([]sdk.FoundBoard, error)
	GetRawDashboardByUID(ctx context.Context, uid string) ([]byte, sdk.BoardProperties, error)
}

type Board struct {
	Title   string
	Content []byte
}

func FetchDashboards(ctx context.Context, c DashboardClient, exportedFolders set.Set[string]) (map[string][]Board, error) {
	foundBoards, err := c.Search(ctx, sdk.SearchType(sdk.SearchTypeDashboard))
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
		if len(exportedFolders) == 0 || !exportedFolders.Contains(board.FolderTitle) {
			slog.Debug("folder not in scope. ignoring", "folderTitle", board.FolderTitle, "title", board.Title)
			continue
		}

		// Get the dashboard JSON model
		rawBoard, _, err := c.GetRawDashboardByUID(ctx, board.UID)
		if err != nil {
			return nil, fmt.Errorf("grafana get board: %w", err)
		}

		boards := append(result[board.FolderTitle], Board{Title: board.Title, Content: rawBoard})
		result[board.FolderTitle] = boards
	}

	return result, nil
}
