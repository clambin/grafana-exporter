package grafana_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/clambin/httpstub"
	"github.com/stretchr/testify/assert"

	"grafana_exporter/internal/grafana"
)

func TestGetDashboardFolders(t *testing.T) {
	var (
		dashboardMap map[string]map[string]string
		folder       map[string]string
		content      string
		ok           bool
		err          error
	)

	dashboardMap, err = grafana.GetAllDashboardsWithHTTPClient(
		"http://example.com",
		"",
		httpstub.NewTestClient(loopback),
	)

	if assert.Nil(t, err) {
		assert.Len(t, dashboardMap, 2)
		folder, ok = dashboardMap["General"]
		assert.True(t, ok)
		assert.Len(t, folder, 1)
		content, ok = folder["db-0-1.json"]
		assert.True(t, ok)
		assert.Equal(t, `"dashboard 2"`, content)
		folder, ok = dashboardMap["folder1"]
		assert.True(t, ok)
		assert.Len(t, folder, 1)
		content, ok = folder["db-1-1.json"]
		assert.True(t, ok)
		assert.Equal(t, `"dashboard 1"`, content)
	}
}

func loopback(req *http.Request) *http.Response {
	switch req.URL.Path {
	case "/api/search":
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(allDashboards)),
		}
	case "/api/dashboards/uid/jQXLLIzRa":
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(dashboard1)),
		}
	case "/api/dashboards/uid/vJMuruVWk":
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(dashboard2)),
		}
	}

	return &http.Response{
		StatusCode: http.StatusNotFound,
	}
}

var (
	allDashboards = `[
  {
    "id": 1,
    "uid": "ReVTKLmRz",
    "title": "folder1",
    "uri": "db/folder1",
    "url": "/dashboards/f/ReVTKLmRz/folder1",
    "slug": "",
    "type": "dash-folder",
    "tags": [],
    "isStarred": false
  },
  {
    "id": 2,
    "uid": "jQXLLIzRa",
    "title": "DB 1.1",
    "uri": "db/db_1_1",
    "url": "/d/jQXLLIzRa/db_1_1",
    "slug": "",
    "type": "dash-db",
    "tags": [],
    "isStarred": false,
    "folderId": 1,
    "folderUid": "ReVTKLmRz",
    "folderTitle": "folder1",
    "folderUrl": "/dashboards/f/ReVTKLmRz/folder1"
  },
  {
    "id": 3,
    "uid": "vJMuruVWk",
    "title": "DB 0.1",
    "uri": "db/db_0_1",
    "url": "/d/vJMuruVWk/DB_0_1",
    "slug": "",
    "type": "dash-db",
    "tags": [],
    "isStarred": false
  }
]`

	dashboard1 = `{ "dashboard": "dashboard 1"}`
	dashboard2 = `{ "dashboard": "dashboard 2"}`
)
