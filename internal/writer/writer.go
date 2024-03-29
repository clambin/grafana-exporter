package writer

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
)

// StorageHandler interface for different storage backends
type StorageHandler interface {
	Initialize() error
	GetCurrent(string) ([]byte, error)
	Add(string, []byte) error
	IsClean() (bool, error)
	Store(msg string) error
}

type Writer struct {
	StorageHandler
	BaseDirectory string
}

func (w *Writer) AddFile(filename string, content []byte) error {
	fullFilename := path.Join(w.BaseDirectory, filename)
	current, err := w.StorageHandler.GetCurrent(fullFilename)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("get current: %w", err)
	} else if bytes.Equal(current, content) {
		return nil
	}

	if err = w.StorageHandler.Add(fullFilename, content); err != nil {
		return fmt.Errorf("add: %w", err)
	}
	return nil
}

func (w *Writer) Store(msg string) error {
	if isClean, err := w.IsClean(); err != nil || isClean {
		return err
	}
	return w.StorageHandler.Store(msg)
}
