package writer_test

import (
	"github.com/clambin/grafana-exporter/internal/writer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"testing"
)

func TestDiskWriter_Write(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "")
	require.NoError(t, err)

	w := writer.NewDiskWriter(tmpdir)

	content := writer.Directories{
		"foo": writer.Files{
			"foo.yaml": []byte("hello"),
		},
		"bar": writer.Files{
			"bar.yaml":   []byte("world"),
			"snafu.yaml": []byte(""),
		},
	}

	require.NoError(t, w.Write(content))

	checkFile(t, tmpdir, "foo/foo.yaml", "hello")
	checkFile(t, tmpdir, "bar/bar.yaml", "world")
	checkFile(t, tmpdir, "bar/snafu.yaml", "")

	assert.NoError(t, w.Write(content))

	require.NoError(t, os.RemoveAll(tmpdir))
}

func checkFile(t *testing.T, tmpdir string, filename string, expected string) {
	t.Helper()

	content, err := os.ReadFile(path.Join(tmpdir, filename))
	require.NoError(t, err)
	assert.Equal(t, expected, string(content))
}

func TestDiskWriter_Write_Error(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "")
	require.NoError(t, err)
	require.NoError(t, os.Chmod(tmpdir, 000))

	w := writer.NewDiskWriter(tmpdir)
	assert.Error(t, w.Write(writer.Directories{
		"foo": writer.Files{
			"foo.yaml": []byte("hello"),
		},
	}))
	require.NoError(t, os.RemoveAll(tmpdir))
}

func TestDiskWriter_Write_Error_2(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "")
	require.NoError(t, err)
	w := writer.NewDiskWriter(tmpdir)

	require.NoError(t, os.MkdirAll(path.Join(tmpdir, "foo"), 000))

	assert.Error(t, w.Write(writer.Directories{
		"foo": writer.Files{
			"foo.yaml": []byte("hello"),
		},
	}))

	require.NoError(t, os.RemoveAll(tmpdir))
}
