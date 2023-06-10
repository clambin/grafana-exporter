package configmap_test

import (
	"github.com/clambin/grafana-exporter/internal/configmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSerialize(t *testing.T) {
	dashboards := map[string][]byte{
		"file-1.json": []byte(`Hello world`),
		"file-2.json": []byte("multi\nline\ncontent"),
	}
	const expected = `kind: ConfigMap
apiVersion: v1
metadata:
  name: foo
  namespace: bar
data:
  file-1.json: Hello world
  file-2.json: |-
    multi
    line
    content
`

	name, output, err := configmap.Serialize(dashboards, "Foo", "bar", "")
	require.NoError(t, err)
	assert.Equal(t, "foo.yml", name)
	assert.Equal(t, expected, string(output))
}

func TestSerialize_WithFolder(t *testing.T) {
	dashboards := map[string][]byte{
		"file-1.json": []byte(`Hello world`),
		"file-2.json": []byte("multi\nline\ncontent"),
	}
	const expected = `kind: ConfigMap
apiVersion: v1
metadata:
  name: foo
  namespace: bar
  labels:
    grafana_dashboard: ""
  annotations:
    grafana_folder: snafu
data:
  file-1.json: Hello world
  file-2.json: |-
    multi
    line
    content
`

	name, output, err := configmap.Serialize(dashboards, "Foo", "bar", "snafu")
	require.NoError(t, err)
	assert.Equal(t, "foo.yml", name)
	assert.Equal(t, expected, string(output))
}
