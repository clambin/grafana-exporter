package grafanatest

import (
	"bytes"
	"github.com/clambin/httpstub"
	"grafana_exporter/internal/grafana"
	"io/ioutil"
	"net/http"
)

// NewWithHTTPClient returns a Grafana Client bound to a stubbed HTTP Server
// Used for unit testing
func NewWithHTTPClient() *grafana.Client {
	return grafana.NewWithHTTPClient(
		"http://example.com",
		"",
		httpstub.NewTestClient(loopback),
	)
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
	case "/api/datasources":
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(allDatasources)),
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

	allDatasources = `[
  {
    "id": 5,
    "orgId": 1,
    "name": "foo",
    "type": "grafana-simple-json-datasource",
    "typeLogoUrl": "public/plugins/grafana-simple-json-datasource/img/simpleJson_logo.svg",
    "access": "proxy",
    "url": "http://example.com:5000",
    "password": "",
    "user": "",
    "database": "",
    "basicAuth": false,
    "isDefault": false,
    "jsonData": {},
    "readOnly": true
  },
  {
    "id": 3,
    "orgId": 1,
    "name": "bar",
    "type": "postgres",
    "typeLogoUrl": "public/app/plugins/datasource/postgres/img/postgresql_logo.svg",
    "access": "proxy",
    "url": "http://postgres.default:5432",
    "password": "",
    "user": "grafana",
    "database": "bar",
    "basicAuth": false,
    "isDefault": false,
    "jsonData": {
      "postgresVersion": 1200,
      "sslmode": "disable"
    },
    "readOnly": true
  },
  {
    "id": 7,
    "orgId": 1,
    "name": "Prometheus",
    "type": "prometheus",
    "typeLogoUrl": "public/app/plugins/datasource/prometheus/img/prometheus_logo.svg",
    "access": "proxy",
    "url": "http://monitoring-prometheus-server.monitoring.svc:80",
    "password": "",
    "user": "",
    "database": "",
    "basicAuth": false,
    "isDefault": true,
    "jsonData": {},
    "readOnly": true
  }
]`
)
