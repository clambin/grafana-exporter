package export_test

import (
	"github.com/clambin/grafana-exporter/export"
	writerMock "github.com/clambin/grafana-exporter/writer/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDashboards_Direct(t *testing.T) {
	writer := &writerMock.Writer{}

	err := export.DashboardProvisioning(writer, true, "")
	assert.NoError(t, err)

	contents, ok := writer.GetFile(".", "dashboards.yml")
	assert.True(t, ok)
	assert.Contains(t, contents, "providers:\n- name: 'dashboards'\n")
}

func TestDashboards_K8S(t *testing.T) {
	writer := &writerMock.Writer{}

	err := export.DashboardProvisioning(writer, false, "monitoring")
	assert.NoError(t, err)

	contents, ok := writer.GetFile(".", "grafana-provisioning-dashboards.yml")
	assert.True(t, ok)
	assert.Contains(t, contents, "providers:\n    - name: 'dashboards'\n")
}
