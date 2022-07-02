package writer_test

import (
	"github.com/clambin/grafana-exporter/writer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"testing"
)

func TestRealWriter_WriteFile(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "")
	require.NoError(t, err)
	w := writer.NewDiskWriter(tmpdir)

	err = w.WriteFile(".", "foo", "bar")
	require.NoError(t, err)
	var content []byte
	content, err = os.ReadFile(path.Join(tmpdir, "foo"))
	require.NoError(t, err)
	assert.Equal(t, "bar", string(content))

	err = w.WriteFile("subdir", "foo", "bar")
	require.NoError(t, err)
	content, err = os.ReadFile(path.Join(tmpdir, "subdir", "foo"))
	require.NoError(t, err)
	assert.Equal(t, "bar", string(content))

	_ = os.RemoveAll(tmpdir)
}

func TestRealWriter_WriteFiles(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "")
	require.NoError(t, err)
	w := writer.NewDiskWriter(tmpdir)

	files := map[string]string{
		"foo": "bar",
		"abc": "xyz",
	}

	err = w.WriteFiles(".", files)
	require.NoError(t, err)

	var content []byte
	content, err = os.ReadFile(path.Join(tmpdir, "foo"))
	require.NoError(t, err)
	assert.Equal(t, "bar", string(content))

	content, err = os.ReadFile(path.Join(tmpdir, "abc"))
	require.NoError(t, err)
	assert.Equal(t, "xyz", string(content))

	err = w.WriteFiles("subdir", files)
	require.NoError(t, err)

	content, err = os.ReadFile(path.Join(tmpdir, "subdir", "foo"))
	require.NoError(t, err)
	assert.Equal(t, "bar", string(content))

	content, err = os.ReadFile(path.Join(tmpdir, "subdir", "abc"))
	require.NoError(t, err)
	assert.Equal(t, "xyz", string(content))

	_ = os.RemoveAll(tmpdir)

}
