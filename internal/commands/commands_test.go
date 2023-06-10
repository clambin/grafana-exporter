package commands_test

import (
	"context"
	"flag"
	"fmt"
	"github.com/clambin/grafana-exporter/internal/fetcher"
	"github.com/clambin/grafana-exporter/internal/writer"
	"github.com/grafana-tools/sdk"
)

var update = flag.Bool("update", false, "generate golden images")
var _ writer.Writer = &fakeWriter{}

type fakeWriter writer.Directories

func (f *fakeWriter) Write(directories writer.Directories) error {
	*f = fakeWriter(directories)
	return nil
}

var _ fetcher.DashboardClient = &fakeDashboardClient{}

type fakeDashboardClient struct{}

func (f fakeDashboardClient) Search(_ context.Context, _ ...sdk.SearchParam) ([]sdk.FoundBoard, error) {
	return []sdk.FoundBoard{
		{Title: "foo", Type: "dash-db", FolderTitle: "bar", UID: "1"},
		{Title: "snafu", Type: "dash-db", FolderTitle: "foobar", UID: "2"},
	}, nil
}

func (f fakeDashboardClient) GetRawDashboardByUID(_ context.Context, uid string) ([]byte, sdk.BoardProperties, error) {
	type dashboardAttribs struct {
		content []byte
		props   sdk.BoardProperties
	}
	var dashboards = map[string]dashboardAttribs{
		"1": {
			content: []byte(`{ "folder": "bar", "title": "foo" }`),
			props: sdk.BoardProperties{
				Type:        "dash-db",
				Slug:        "foo",
				FolderTitle: "bar",
			},
		},
		"2": {
			content: []byte(`{ "folder": "foobar", "title": "snafu" }`),
			props: sdk.BoardProperties{
				Type:        "dash-db",
				Slug:        "foo",
				FolderTitle: "bar",
			},
		},
	}
	var err error
	attribs, ok := dashboards[uid]
	if !ok {
		err = fmt.Errorf("invalid uid: %s", uid)
	}
	return attribs.content, attribs.props, err
}
