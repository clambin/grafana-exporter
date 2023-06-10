package fetcher

import (
	gapi "github.com/grafana/grafana-api-golang-client"
)

type DataSourcesClient interface {
	DataSources() ([]*gapi.DataSource, error)
}
