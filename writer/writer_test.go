package writer_test

import (
	"github.com/clambin/grafana-exporter/writer"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestRealWriter_WriteFile(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "")
	if assert.NoError(t, err) == false {
		t.Fatal()
	}

	w := writer.NewDiskWriter(tmpdir)

	err = w.WriteFile(".", "foo", "bar")
	assert.NoError(t, err)
	var content []byte
	content, err = os.ReadFile(path.Join(tmpdir, "foo"))
	assert.NoError(t, err)
	assert.Equal(t, "bar", string(content))

	err = w.WriteFile("subdir", "foo", "bar")
	assert.NoError(t, err)
	content, err = os.ReadFile(path.Join(tmpdir, "subdir", "foo"))
	assert.NoError(t, err)
	assert.Equal(t, "bar", string(content))

	_ = os.RemoveAll(tmpdir)
}

func TestRealWriter_WriteFiles(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "")
	if assert.NoError(t, err) == false {
		t.Fatal()
	}
	w := writer.NewDiskWriter(tmpdir)

	files := map[string]string{
		"foo": "bar",
		"abc": "xyz",
	}

	err = w.WriteFiles(".", files)
	assert.NoError(t, err)

	var content []byte
	content, err = os.ReadFile(path.Join(tmpdir, "foo"))
	assert.NoError(t, err)
	assert.Equal(t, "bar", string(content))

	content, err = os.ReadFile(path.Join(tmpdir, "abc"))
	assert.NoError(t, err)
	assert.Equal(t, "xyz", string(content))

	err = w.WriteFiles("subdir", files)
	assert.NoError(t, err)

	content, err = os.ReadFile(path.Join(tmpdir, "subdir", "foo"))
	assert.NoError(t, err)
	assert.Equal(t, "bar", string(content))

	content, err = os.ReadFile(path.Join(tmpdir, "subdir", "abc"))
	assert.NoError(t, err)
	assert.Equal(t, "xyz", string(content))

	_ = os.RemoveAll(tmpdir)

}
