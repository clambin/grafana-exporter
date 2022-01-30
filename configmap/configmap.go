package configmap

import (
	"bytes"
	"github.com/gosimple/slug"
	"gopkg.in/yaml.v3"
)

// Serialize creates a ConfigMap structure and serializes it into a byte slice
// so we can store it in a yaml file
func Serialize(name, namespace, folder string, files map[string]string) (string, string, error) {
	type metadata struct {
		Name        string            `yaml:"name"`
		Namespace   string            `yaml:"namespace"`
		Labels      map[string]string `yaml:"labels,omitempty"`
		Annotations map[string]string `yaml:"annotations,omitempty"`
	}

	type configMap struct {
		Kind       string            `yaml:"kind"`
		APIVersion string            `yaml:"apiVersion"`
		Metadata   metadata          `yaml:"metadata"`
		Data       map[string]string `yaml:"data"`
	}

	mapName := slug.Make(name)
	configmap := configMap{
		Kind:       "ConfigMap",
		APIVersion: "v1",
		Metadata:   metadata{Name: mapName, Namespace: namespace},
		Data:       files,
	}

	if folder != "" {
		configmap.Metadata.Labels = map[string]string{"grafana_dashboard": ""}
		configmap.Metadata.Annotations = map[string]string{"grafana_folder": folder}
	}

	var b bytes.Buffer
	encoder := yaml.NewEncoder(&b)
	encoder.SetIndent(2)
	err := encoder.Encode(configmap)
	return mapName + ".yml", string(b.Bytes()), err
}
