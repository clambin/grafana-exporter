package grafana_test

import (
	"bytes"
	"github.com/clambin/httpstub"
	"github.com/stretchr/testify/assert"
	"grafana_exporter/internal/grafana"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestGetDatasources(t *testing.T) {
	datasourceMap, err := grafana.GetDatasourcesWithHTTPClient("http://example.com", "",
		httpstub.NewTestClient(loopbackDatasources),
	)

	var ok bool
	assert.Nil(t, err)
	assert.Len(t, datasourceMap, 1)
	content, ok := datasourceMap["datasources.yml"]
	assert.True(t, ok)
	assert.Equal(t, expected, content)
}

func loopbackDatasources(req *http.Request) *http.Response {
	switch req.URL.Path {
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

const (
	expected = `apiVersion: 1
datasources:
  - id: 5
    orgid: 1
    name: foo
    type: grafana-simple-json-datasource
    access: proxy
    url: http://example.com:5000
    password: ""
    user: ""
    database: ""
    basicauth: false
    basicauthuser: null
    basicauthpassword: null
    isdefault: false
    jsondata: {}
    securejsondata: null
  - id: 3
    orgid: 1
    name: bar
    type: postgres
    access: proxy
    url: http://postgres.default:5432
    password: ""
    user: grafana
    database: bar
    basicauth: false
    basicauthuser: null
    basicauthpassword: null
    isdefault: false
    jsondata:
        postgresVersion: 1200
        sslmode: disable
    securejsondata: null
  - id: 7
    orgid: 1
    name: Prometheus
    type: prometheus
    access: proxy
    url: http://monitoring-prometheus-server.monitoring.svc:80
    password: ""
    user: ""
    database: ""
    basicauth: false
    basicauthuser: null
    basicauthpassword: null
    isdefault: true
    jsondata: {}
    securejsondata: null
`

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
