package exporter_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"grafana_exporter/internal/exporter"
	"os"
	"testing"
)

func TestExporter(t *testing.T) {
	log := newLogger()
	dir := os.TempDir()
	err := exporter.NewWithLogger(
		"http://grafana.192.168.0.11.nip.io",
		`admin:catch22`,
		dir,
		"monitoring",
		log.writeFile,
	).Export()
	assert.Nil(t, err)

	fmt.Println(dir)
}

type logger struct {
	output map[string]map[string][]byte
}

func newLogger() *logger {
	return &logger{
		output: make(map[string]map[string][]byte),
	}
}

func (log *logger) writeFile(directory, filename string, content []byte) {
	var ok bool

	if _, ok = log.output[directory]; ok == false {
		log.output[directory] = make(map[string][]byte)
	}
	log.output[directory][filename] = content
}
