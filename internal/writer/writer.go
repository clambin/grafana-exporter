package writer

import (
	"errors"
	"fmt"
	"os"
	"path"
)

type Files map[string][]byte
type Directories map[string]Files

type Writer interface {
	Write(Directories) error
}

var _ Writer = DiskWriter{}

type DiskWriter struct {
	baseDirectory string
}

func NewDiskWriter(baseDirectory string) DiskWriter {
	return DiskWriter{baseDirectory: baseDirectory}
}

func (dw DiskWriter) Write(directories Directories) error {
	for directory, files := range directories {
		if err := os.Mkdir(path.Join(dw.baseDirectory, directory), 0755); err != nil && !errors.Is(err, os.ErrExist) {
			return fmt.Errorf("create directory %s: %w", directory, err)
		}

		for name, content := range files {
			if err := os.WriteFile(path.Join(dw.baseDirectory, directory, name), content, 0644); err != nil {
				return fmt.Errorf("write %s: %w", path.Join(directory, name), err)
			}
		}
	}
	return nil
}
