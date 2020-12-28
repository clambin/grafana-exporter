package grafana

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gosimple/slug"
	"github.com/grafana-tools/sdk"
	"gopkg.in/yaml.v3"
	"net/http"
)

// GetAllDashboards retrieves all dashboards in Grafana.
// Dashboards may be located in folders.  GetAllDashboards therefore returns
// a map of folders, each of which holds a map of dashboards with their filename
// Names are converted to slugs for use in clusters and/or file systems.
//
// E.g. if Grafana has
//     Folder 1:  Dashboard 1
//     Folder 2:  Dashboard 2, Dashboard 3
// then this function returns
// map
// +-> folder-1
// |"  +-> dashboard-1.json -> json model of dashboard 1
// +-> folder-2
//     +-> dashboard-2.json -> json model of dashboard 2
//     +-> dashboard-3.json -> json model of dashboard 3
func GetAllDashboards(url, apiToken string) (map[string]map[string]string, error) {
	return GetAllDashboardsWithHTTPClient(url, apiToken, sdk.DefaultHTTPClient)
}

// GetAllDashboardsWithHTTPClient uses the specified http.Client to call the Grafana APIs
// Used in unit testing to stub the API calls
func GetAllDashboardsWithHTTPClient(url, apiToken string, client *http.Client) (map[string]map[string]string, error) {
	var (
		err         error
		foundBoards []sdk.FoundBoard
		rawBoard    []byte
	)
	result := make(map[string]map[string]string)

	ctx := context.Background()
	c := sdk.NewClient(url, apiToken, client)

	// Get all dashboards
	if foundBoards, err = c.Search(ctx, sdk.SearchType(sdk.SearchTypeDashboard)); err == nil {
		for _, link := range foundBoards {
			// Only process dashboards, not folders
			if link.Type == "dash-db" {
				// Get the dashboard JSON model
				if rawBoard, _, err = c.GetRawDashboardByUID(ctx, link.UID); err == nil {
					// The "General" board has an empty title in Grafana
					if link.FolderTitle == "" {
						link.FolderTitle = "General"
					}
					// First dashboard for this folder? Create the map
					if _, ok := result[link.FolderTitle]; ok == false {
						result[link.FolderTitle] = make(map[string]string)
					}
					// Reformat the JSON stream to store it properly in YAML
					var buffer bytes.Buffer
					_ = json.Indent(&buffer, rawBoard, "", "  ")
					// Store it in the map
					result[link.FolderTitle][slug.Make(link.Title)+".json"] = string(buffer.Bytes())
				} else {
					break
				}
			}
		}
	}

	return result, err
}

// GetDatasources retrieves all datasources in Grafana.
// For simplicity, we'll store these in one config file 'datasources.yml'
// So the returning map will only have one element.
func GetDatasources(url, apiToken string) (map[string]string, error) {
	return GetDatasourcesWithHTTPClient(url, apiToken, sdk.DefaultHTTPClient)
}

// GetDatasourcesWithHTTPClient uses the specified http.Client to call the Grafana APIs
// Used in unit testing to stub the API calls
func GetDatasourcesWithHTTPClient(url, apiToken string, client *http.Client) (map[string]string, error) {
	var (
		err         error
		datasources []sdk.Datasource
		dsPacked    []byte
	)
	result := make(map[string]string)
	ctx := context.Background()
	c := sdk.NewClient(url, apiToken, client)

	if datasources, err = c.GetAllDatasources(ctx); err == nil {
		// datasource provisioning uses apiVersion / datasources layout
		type dataSource struct {
			APIVersion  int              `yaml:"apiVersion"`
			Datasources []sdk.Datasource `yaml:"datasources"`
		}
		var wrapper = dataSource{1, datasources}

		if dsPacked, err = yaml.Marshal(&wrapper); err == nil {
			result["datasources.yml"] = string(dsPacked)
		}
	}
	return result, err
}
