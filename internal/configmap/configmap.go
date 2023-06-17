package configmap

import (
	"bytes"
	"github.com/gosimple/slug"
	"gopkg.in/yaml.v3"
)

type configMap struct {
	Kind       string            `yaml:"kind"`
	APIVersion string            `yaml:"apiVersion"`
	Metadata   metadata          `yaml:"metadata"`
	Data       map[string]string `yaml:"data"`
}

type metadata struct {
	Name        string            `yaml:"name"`
	Namespace   string            `yaml:"namespace"`
	Labels      map[string]string `yaml:"labels,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

// Serialize creates a ConfigMap that contains the specified files. If the files are dashboards (folder is not blank),
// it will add the necessary metadata (i.e. grafana_dashboard label and grafana_folder annotation, so Grafana will
// detect them as dashboards and import them.
func Serialize(files map[string][]byte, configMapName, namespace, folder string) (string, []byte, error) {
	filesAsString := make(map[string]string)
	for filename, content := range files {
		filesAsString[filename] = string(content)
	}

	md := metadata{Name: slug.Make(configMapName), Namespace: namespace}
	if folder != "" {
		md.Labels = map[string]string{"grafana_dashboard": ""}
		md.Annotations = map[string]string{"grafana_folder": folder}
	}

	configmap := configMap{
		Kind:       "ConfigMap",
		APIVersion: "v1",
		Metadata:   md,
		Data:       filesAsString,
	}

	var b bytes.Buffer
	encoder := yaml.NewEncoder(&b)
	encoder.SetIndent(2)
	err := encoder.Encode(configmap)
	_ = encoder.Close()
	return configmap.Metadata.Name + ".yml", b.Bytes(), err
}
