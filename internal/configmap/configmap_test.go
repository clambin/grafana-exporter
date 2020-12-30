package configmap_test

import (
	"github.com/stretchr/testify/assert"
	"grafana_exporter/internal/configmap"
	"testing"
)

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

func TestSerialize(t *testing.T) {
	dashboards := map[string]string{
		"file-1.json": `Hello world`,
		"file-2.json": "multi\nline\ncontent",
	}

	name, output, err := configmap.Serialize("Foo", "bar", dashboards)
	assert.Nil(t, err)
	assert.Equal(t, "foo.yml", name)
	assert.Equal(t, expected, string(output))
}
