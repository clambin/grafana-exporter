package export_test

import (
	"github.com/clambin/grafana-exporter/export"
	"github.com/clambin/grafana-exporter/grafana"
	grafanaMock "github.com/clambin/grafana-exporter/grafana/mock"
	writerMock "github.com/clambin/grafana-exporter/writer/mock"
	"github.com/gosimple/slug"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestDashboards(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(grafanaMock.ServerHandler))
	defer server.Close()
	client := grafana.New(server.URL, "")

	testCases := []struct {
		name    string
		direct  bool
		folders []string
		pass    bool
		files   map[string][]string
	}{
		{
			name:   "k8s",
			direct: false,
			pass:   true,
			files:  map[string][]string{".": {"grafana-dashboards-general.yml", "grafana-dashboards-folder1.yml"}},
		},
		{
			name:   "direct",
			direct: true,
			pass:   true,
			files: map[string][]string{
				"folder1": {"db-1-1.json"},
				"General": {"db-0-1.json"},
			},
		},
		{
			name:    "filtered",
			direct:  false,
			pass:    true,
			folders: []string{"folder1"},
			files:   map[string][]string{".": {"grafana-dashboards-folder1.yml"}},
		},
	}

	for _, tt := range testCases {
		w := writerMock.Writer{}

		if err := export.Dashboards(client, &w, tt.direct, "namespace", tt.folders); !tt.pass {
			assert.Error(t, err, tt.name)
			continue
		} else {
			require.NoError(t, err, tt.name)
		}

		for dir, files := range tt.files {
			count, found := w.Count(dir)
			require.True(t, found, tt.name+": "+dir)
			assert.Equal(t, len(tt.files[dir]), count)

			for _, file := range files {
				var content string
				content, found = w.GetFile(dir, file)
				require.True(t, found, tt.name+": "+dir+"/"+file)

				assert.NotEmpty(t, content, tt.name)

				gp := filepath.Join("testdata", t.Name()+"_"+tt.name+"_"+slug.Make(dir)+"_"+slug.Make(file)+".golden")
				if *update {
					require.NoError(t, os.WriteFile(gp, []byte(content), 0644))
				}
				golden, err := os.ReadFile(gp)
				require.NoError(t, err)
				assert.Equal(t, string(golden), content)
			}
		}
	}
}
