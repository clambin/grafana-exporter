package fetcher_test

import (
	"fmt"
	"github.com/clambin/go-common/set"
	"github.com/clambin/grafana-exporter/internal/fetcher"
	gapi "github.com/grafana/grafana-api-golang-client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFetchDashboards(t *testing.T) {
	testcases := []struct {
		name    string
		folders set.Set[string]
		boards  []string
		titles  []string
	}{
		{
			name:    "all folders",
			folders: set.New[string](),
			boards:  []string{"foo", "bar"},
			titles:  []string{"board 1", "board 2"},
		},
		{
			name:    "filters",
			folders: set.New("foo"),
			boards:  []string{"foo"},
			titles:  []string{"board 1"},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fetcher.FetchDashboards(&fakeDashboardFetcher{}, tt.folders)
			require.NoError(t, err)
			assert.Len(t, result, len(tt.boards))

			for idx, board := range tt.boards {
				b, ok := result[board]
				require.Truef(t, ok, board)
				require.Lenf(t, b, 1, board)
				assert.Equalf(t, tt.titles[idx], b[0].Title, board)
			}
		})
	}
}

var _ fetcher.DashboardClient = &fakeDashboardFetcher{}

type fakeDashboardFetcher struct {
}

func (f fakeDashboardFetcher) Dashboards() ([]gapi.FolderDashboardSearchResponse, error) {
	return []gapi.FolderDashboardSearchResponse{
		{UID: "1", Title: "board 1", Type: "dash-db", FolderTitle: "foo"},
		{UID: "2", Title: "board 2", Type: "dash-db", FolderTitle: "bar"},
		{UID: "3", Title: "foo", Type: "folder", FolderTitle: ""},
		{UID: "4", Title: "bar", Type: "folder", FolderTitle: ""},
	}, nil
}

func (f fakeDashboardFetcher) DashboardByUID(uid string) (*gapi.Dashboard, error) {
	dashboards := map[string]*gapi.Dashboard{
		"1": {Model: map[string]any{"foo": "bar"}},
		"2": {Model: map[string]any{"bar": "foo"}},
	}

	dashboard, ok := dashboards[uid]
	if !ok {
		return nil, fmt.Errorf("invalid dashboard uid: %s", uid)
	}
	return dashboard, nil
}
