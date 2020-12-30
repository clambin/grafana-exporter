package configmap

import (
	"bytes"
	"github.com/gosimple/slug"
	"gopkg.in/yaml.v3"
)

// Serialize creates a ConfigMap structure and serializes it into a byte slice
// so we can store it in a yaml file
func Serialize(name, namespace string, files map[string]string) (string, []byte, error) {
	type metadata struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	}

	type configMap struct {
		Kind       string            `yaml:"kind"`
		APIVersion string            `yaml:"apiVersion"`
		Metadata   metadata          `yaml:"metadata"`
		Data       map[string]string `yaml:"data"`
	}

	var (
		mapName   = slug.Make(name)
		configmap = configMap{
			"ConfigMap", "v1",
			metadata{mapName, namespace},
			files,
		}
		b bytes.Buffer
	)

	encoder := yaml.NewEncoder(&b)
	encoder.SetIndent(2)
	err := encoder.Encode(configmap)
	return mapName + ".yml", b.Bytes(), err
}
