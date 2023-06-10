package commands_test

import (
	"flag"
	"fmt"
	"github.com/clambin/grafana-exporter/internal/fetcher"
	"github.com/clambin/grafana-exporter/internal/writer"
	gapi "github.com/grafana/grafana-api-golang-client"
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

func (f fakeDashboardClient) Dashboards() ([]gapi.FolderDashboardSearchResponse, error) {
	return []gapi.FolderDashboardSearchResponse{
		{Title: "foo", Type: "dash-db", FolderTitle: "bar", UID: "1"},
		{Title: "snafu", Type: "dash-db", FolderTitle: "foobar", UID: "2"},
	}, nil
}

func (f fakeDashboardClient) DashboardByUID(uid string) (*gapi.Dashboard, error) {
	var dashboards = map[string]*gapi.Dashboard{
		"1": {Model: map[string]any{"folder": "bar", "title": "foo"}},
		"2": {Model: map[string]any{"folder": "foobar", "title": "snafu"}},
	}
	if dashboard, ok := dashboards[uid]; ok {
		return dashboard, nil
	}
	return nil, fmt.Errorf("invalid uid: %s", uid)
}
