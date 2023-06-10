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
	result, err := fetcher.FetchDashboards(&fakeDashboardFetcher{}, set.Create("foo"))
	require.NoError(t, err)
	assert.Len(t, result, 1)

	boards, ok := result["foo"]
	require.True(t, ok)
	require.Len(t, boards, 1)

	assert.Equal(t, "board 1", boards[0].Title)
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
