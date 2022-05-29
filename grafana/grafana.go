package grafana

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
func (client *Client) GetAllDashboards(ctx context.Context, exportedFolders []string) (result map[string]map[string]string, err error) {
	c := sdk.NewClient(client.url, client.apiToken, client.apiClient)

	// Get all dashboards
	var foundBoards []sdk.FoundBoard
	foundBoards, err = c.Search(ctx, sdk.SearchType(sdk.SearchTypeDashboard))

	if err != nil {
		err = fmt.Errorf("failed to get dashboard from grafana: %w", err)
		return
	}

	result = make(map[string]map[string]string)
	for _, link := range foundBoards {
		log.WithFields(log.Fields{
			"title":  link.Title,
			"type":   link.Type,
			"folder": link.FolderTitle,
		}).Debug("dashboard found")

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
		if !validFolder(link.FolderTitle, exportedFolders) {
			log.WithField("folderTitle", link.FolderTitle).Debug("folder not in scope. ignoring")
			continue
		}

		// Get the dashboard JSON model
		var rawBoard []byte
		rawBoard, _, err = c.GetRawDashboardByUID(ctx, link.UID)
		if err != nil {
			log.Warnf("failed to get dashboard %s: %s", link.Title, err.Error())
			continue
		}

		// Reformat the JSON stream to store it properly in YAML
		var buffer bytes.Buffer
		_ = json.Indent(&buffer, rawBoard, "", "  ")

		// First dashboard for this folder? Create the map
		if _, ok := result[link.FolderTitle]; !ok {
			result[link.FolderTitle] = make(map[string]string)
		}

		// Store it in the map
		result[link.FolderTitle][slug.Make(link.Title)+".json"] = buffer.String()
		log.Debug("Stored")
	}
	return
}

// GetDataSources retrieves all dataSources in Grafana.
// For simplicity, we'll store these in one config file 'datasources.yml'
// So the returning map will only have one element.
func (client *Client) GetDataSources(ctx context.Context) (result map[string]string, err error) {
	c := sdk.NewClient(client.url, client.apiToken, client.apiClient)

	var dataSources []sdk.Datasource
	if dataSources, err = c.GetAllDatasources(ctx); err != nil {
		return
	}

	wrapper := struct {
		APIVersion  int              `yaml:"apiVersion"`
		DataSources []sdk.Datasource `yaml:"datasources"`
	}{
		APIVersion:  1,
		DataSources: dataSources,
	}

	var buffer bytes.Buffer
	encoder := yaml.NewEncoder(&buffer)
	encoder.SetIndent(2)

	err = encoder.Encode(&wrapper)
	_ = encoder.Close()

	if err == nil {
		result = map[string]string{"datasources.yml": buffer.String()}
	}

	return
}

func validFolder(folder string, folders []string) bool {
	if len(folders) == 0 {
		return true
	}
	for _, f := range folders {
		if f == folder {
			return true
		}
	}
	return false
}
