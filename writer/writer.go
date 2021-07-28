package writer

import (
	"os"
	"path"
)

// Writer implements an interface to save one or more files
type Writer interface {
	WriteFiles(directory string, files map[string]string) (err error)
	WriteFile(directory, filename string, content string) (err error)
}

// DiskWriter implements the Writer interface to save files to disk
type DiskWriter struct {
	Directory string
}

// NewDiskWriter creates a new DiskWriter
func NewDiskWriter(directory string) *DiskWriter {
	return &DiskWriter{
		Directory: directory,
	}
}

// WriteFiles saves files to the specified directory
func (w *DiskWriter) WriteFiles(directory string, files map[string]string) (err error) {
	for fileName, fileContents := range files {
		err = w.WriteFile(directory, fileName, fileContents)

		if err != nil {
			break
		}
	}

	return
}

// WriteFile saves one file to the specified directory
func (w *DiskWriter) WriteFile(directory, filename string, content string) (err error) {
	targetDir := path.Join(w.Directory, directory)

	err = os.MkdirAll(targetDir, 0755)

	if err == nil {
		err = os.WriteFile(path.Join(targetDir, filename), []byte(content), 0644)
	}

	return
}
