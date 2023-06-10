package fetcher

import (
	"context"
	"github.com/grafana-tools/sdk"
)

type DataSourcesClient interface {
	GetAllDatasources(ctx context.Context) ([]sdk.Datasource, error)
}

func FetchDataSources(ctx context.Context, c DataSourcesClient) ([]sdk.Datasource, error) {
	return c.GetAllDatasources(ctx)
}
