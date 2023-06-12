package writer_test

import (
	"github.com/clambin/grafana-exporter/internal/writer"
	"github.com/clambin/grafana-exporter/internal/writer/fs"
	"github.com/clambin/grafana-exporter/internal/writer/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"testing"
)

func TestMockedWriter(t *testing.T) {
	type mockParameters struct {
		methodName   string
		arguments    []any
		returnValues []any
	}
	testcases := []struct {
		name           string
		mockParameters []mockParameters
		wantWriteErr   assert.ErrorAssertionFunc
		wantFlushErr   assert.ErrorAssertionFunc
	}{
		{
			name: "change",
			mockParameters: []mockParameters{
				{methodName: "Initialize", returnValues: []any{nil}},
				{methodName: "GetCurrent", arguments: []any{"/tmp/foo/bar.txt"}, returnValues: []any{[]byte("old"), nil}},
				{methodName: "Add", arguments: []any{"/tmp/foo/bar.txt", []byte("hello")}, returnValues: []any{nil}},
				{methodName: "IsClean", returnValues: []any{false, nil}},
				{methodName: "Store", returnValues: []any{nil}},
			},
			wantWriteErr: assert.NoError,
			wantFlushErr: assert.NoError,
		},
		{
			name: "no change",
			mockParameters: []mockParameters{
				{methodName: "Initialize", returnValues: []any{nil}},
				{methodName: "GetCurrent", arguments: []any{"/tmp/foo/bar.txt"}, returnValues: []any{[]byte("hello"), nil}},
				{methodName: "IsClean", returnValues: []any{true, nil}},
			},
			wantWriteErr: assert.NoError,
			wantFlushErr: assert.NoError,
		},
		{
			name: "new",
			mockParameters: []mockParameters{
				{methodName: "Initialize", returnValues: []any{nil}},
				{methodName: "GetCurrent", arguments: []any{"/tmp/foo/bar.txt"}, returnValues: []any{nil, os.ErrNotExist}},
				{methodName: "Add", arguments: []any{"/tmp/foo/bar.txt", []byte("hello")}, returnValues: []any{nil}},
				{methodName: "IsClean", returnValues: []any{false, nil}},
				{methodName: "Store", returnValues: []any{nil}},
			},
			wantWriteErr: assert.NoError,
			wantFlushErr: assert.NoError,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			r := mocks.NewStorageHandler(t)
			w := writer.Writer{StorageHandler: r, BaseDirectory: "/tmp"}

			for _, m := range tt.mockParameters {
				r.On(m.methodName, m.arguments...).Return(m.returnValues...).Once()
			}

			require.NoError(t, r.Initialize())
			tt.wantWriteErr(t, w.AddFile("foo/bar.txt", []byte("hello")))
			tt.wantFlushErr(t, w.Store())
		})
	}
}

func TestFSWriter(t *testing.T) {
	testcases := []struct {
		name    string
		content []byte
		clean   bool
	}{
		{
			name:    "new",
			content: []byte("hello"),
		},
		{
			name:    "change",
			content: []byte("world"),
		},
		{
			name:    "no change",
			content: []byte("world"),
			clean:   true,
		},
	}

	tmpdir, err := os.MkdirTemp("", "")
	require.NoError(t, err)

	w := writer.Writer{
		StorageHandler: &fs.Client{},
		BaseDirectory:  tmpdir,
	}
	require.NoError(t, w.Initialize())

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			assert.NoError(t, w.AddFile("foo/bar.txt", tt.content))

			clean, err := w.IsClean()
			require.NoError(t, err)
			assert.Equal(t, tt.clean, clean)

			assert.NoError(t, w.Store())

			content, err := os.ReadFile(path.Join(tmpdir, "foo/bar.txt"))
			require.NoError(t, err)
			assert.Equal(t, tt.content, content)
		})
	}

	assert.NoError(t, os.RemoveAll(tmpdir))
}
