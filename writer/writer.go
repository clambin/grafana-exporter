package writer

import (
	"os"
	"path"
)

type Writer interface {
	WriteFiles(directory string, files map[string]string) (err error)
	WriteFile(directory, filename string, content string) (err error)
}

type RealWriter struct {
	Directory string
}

func NewWriter(directory string) *RealWriter {
	return &RealWriter{
		Directory: directory,
	}
}

func (w *RealWriter) WriteFiles(directory string, files map[string]string) (err error) {
	for fileName, fileContents := range files {
		err = w.WriteFile(directory, fileName, fileContents)

		if err != nil {
			break
		}
	}

	return
}

func (w *RealWriter) WriteFile(directory, filename string, content string) (err error) {
	targetDir := path.Join(w.Directory, directory)

	err = os.MkdirAll(targetDir, 0755)

	if err == nil {
		err = os.WriteFile(path.Join(targetDir, filename), []byte(content), 0644)
	}

	return
}
