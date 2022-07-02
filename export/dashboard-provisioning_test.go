package export_test

import (
	"flag"
	"github.com/clambin/grafana-exporter/export"
	writerMock "github.com/clambin/grafana-exporter/writer/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

var update = flag.Bool("update", false, "update .golden files")

func TestDashboardProvisioning_Direct(t *testing.T) {
	writer := &writerMock.Writer{}
	err := export.DashboardProvisioning(writer, true, "")
	assert.NoError(t, err)
	contents, ok := writer.GetFile(".", "dashboards.yml")
	assert.True(t, ok)

	gp := filepath.Join("testdata", t.Name()+".golden")
	if *update {
		err = os.WriteFile(gp, []byte(contents), 0644)
		require.NoError(t, err)
	}

	var golden []byte
	golden, err = os.ReadFile(gp)
	require.NoError(t, err)
	assert.Equal(t, string(golden), contents)
}

func TestDashboardProvisioning_K8S(t *testing.T) {
	writer := &writerMock.Writer{}
	err := export.DashboardProvisioning(writer, false, "monitoring")
	assert.NoError(t, err)
	contents, ok := writer.GetFile(".", "grafana-provisioning-dashboards.yml")
	assert.True(t, ok)

	gp := filepath.Join("testdata", t.Name()+".golden")
	if *update {
		err = os.WriteFile(gp, []byte(contents), 0644)
		require.NoError(t, err)
	}

	var golden []byte
	golden, err = os.ReadFile(gp)
	require.NoError(t, err)
	assert.Equal(t, string(golden), contents)
}
