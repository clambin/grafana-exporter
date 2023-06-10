package fetcher_test

import (
	"context"
	"errors"
	"github.com/clambin/go-common/set"
	"github.com/clambin/grafana-exporter/internal/fetcher"
	"github.com/grafana-tools/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFetchDashboards(t *testing.T) {
	result, err := fetcher.FetchDashboards(context.Background(), &fakeDashboardFetcher{}, set.Create("foo"))
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

func (f fakeDashboardFetcher) Search(_ context.Context, _ ...sdk.SearchParam) ([]sdk.FoundBoard, error) {
	return []sdk.FoundBoard{
		{UID: "1", Title: "board 1", Type: "dash-db", FolderTitle: "foo"},
		{UID: "2", Title: "board 2", Type: "dash-db", FolderTitle: "bar"},
		{UID: "3", Title: "foo", Type: "folder", FolderTitle: ""},
		{UID: "4", Title: "bar", Type: "folder", FolderTitle: ""},
	}, nil
}

func (f fakeDashboardFetcher) GetRawDashboardByUID(_ context.Context, uid string) ([]byte, sdk.BoardProperties, error) {
	var props sdk.BoardProperties
	if uid != "1" {
		return nil, props, errors.New("not found")
	}
	return []byte(`{
"foo": "bar",
"bar": "foo"
`), props, nil
}
