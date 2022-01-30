package configmap_test

import (
	configmap2 "github.com/clambin/grafana-exporter/configmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSerialize(t *testing.T) {
	dashboards := map[string]string{
		"file-1.json": `Hello world`,
		"file-2.json": "multi\nline\ncontent",
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

	name, output, err := configmap2.Serialize("Foo", "bar", "", dashboards)
	require.NoError(t, err)
	assert.Equal(t, "foo.yml", name)
	assert.Equal(t, expected, output)
}

func TestSerialize_WithFolder(t *testing.T) {
	dashboards := map[string]string{
		"file-1.json": `Hello world`,
		"file-2.json": "multi\nline\ncontent",
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

	name, output, err := configmap2.Serialize("Foo", "bar", "snafu", dashboards)
	require.NoError(t, err)
	assert.Equal(t, "foo.yml", name)
	assert.Equal(t, expected, output)
}
