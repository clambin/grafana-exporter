package grafana

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gosimple/slug"
	"github.com/grafana-tools/sdk"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"net/http"
)

// Client to call Grafana APIs
type Client struct {
	apiClient *http.Client
	url       string
	apiToken  string
}

// New creates a new Client
func New(url, apiToken string) *Client {
	return NewWithHTTPClient(url, apiToken, sdk.DefaultHTTPClient)
}

// NewWithHTTPClient creates a Client with a specified http.Client
// Used to stub API calls during unit testing
func NewWithHTTPClient(url, apiToken string, httpClient *http.Client) *Client {
	return &Client{
		url:       url,
		apiToken:  apiToken,
		apiClient: httpClient,
	}
}

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
func (client *Client) GetAllDashboards(ctx context.Context, exportedFolders []string) (map[string]map[string]string, error) {
	var (
		err         error
		foundBoards []sdk.FoundBoard
		rawBoard    []byte
	)
	result := make(map[string]map[string]string)

	c := sdk.NewClient(client.url, client.apiToken, client.apiClient)

	// Get all dashboards
	if foundBoards, err = c.Search(ctx, sdk.SearchType(sdk.SearchTypeDashboard)); err == nil {
		for _, link := range foundBoards {
			log.Debugf("considering %s (type: %s; folder: %s)", link.Title, link.Type, link.FolderTitle)
			// Only process dashboards, not folders
			if link.Type != "dash-db" {
				log.WithField("type", link.Type).Debug("wrong type. ignoring")
				continue
			}
			// Het kind moet toch een naam hebben
			if link.FolderTitle == "" {
				link.FolderTitle = "General"
			}
			// Only export if the dashboard is in a specified folder
			if len(exportedFolders) > 0 && validFolder(link.FolderTitle, exportedFolders) == false {
				log.WithField("folderTitle", link.FolderTitle).Debug("folder not in scope. ignoring")
				continue
			}
			// Get the dashboard JSON model
			if rawBoard, _, err = c.GetRawDashboardByUID(ctx, link.UID); err == nil {
				// Reformat the JSON stream to store it properly in YAML
				var buffer bytes.Buffer
				_ = json.Indent(&buffer, rawBoard, "", "  ")
				// First dashboard for this folder? Create the map
				if _, ok := result[link.FolderTitle]; ok == false {
					result[link.FolderTitle] = make(map[string]string)
				}
				// Store it in the map
				result[link.FolderTitle][slug.Make(link.Title)+".json"] = string(buffer.Bytes())
				log.Debug("Stored")
			} else {
				log.Warnf("failed to get dashboard %s: %s", link.Title, err.Error())
			}
		}
	}

	return result, err
}

// GetDataSources retrieves all dataSources in Grafana.
// For simplicity, we'll store these in one config file 'datasources.yml'
// So the returning map will only have one element.
func (client *Client) GetDataSources(ctx context.Context) (map[string]string, error) {
	var (
		err         error
		datasources []sdk.Datasource
		dsPacked    []byte
	)
	result := make(map[string]string)
	c := sdk.NewClient(client.url, client.apiToken, client.apiClient)

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

func validFolder(folder string, folders []string) bool {
	for _, f := range folders {
		if f == folder {
			return true
		}
	}
	return false
}
