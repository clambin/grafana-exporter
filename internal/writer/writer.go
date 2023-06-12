package writer

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
)

// StorageHandler interface for different storage backends
//
//go:generate mockery --name StorageHandler
type StorageHandler interface {
	Initialize() error
	GetCurrent(string) ([]byte, error)
	Mkdir(string) error
	Add(string, []byte) error
	IsClean() (bool, error)
	Store() error
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

	dirname := path.Dir(fullFilename)
	if err := os.MkdirAll(dirname, 0755); err != nil && !errors.Is(err, os.ErrExist) {
		return fmt.Errorf("mkdir: %w", err)
	}

	if err := w.StorageHandler.Add(fullFilename, content); err != nil {
		return fmt.Errorf("add: %w", err)
	}
	return nil
}

func (w *Writer) Store() error {
	if isClean, err := w.IsClean(); err != nil || isClean {
		return err
	}
	return w.StorageHandler.Store()
}
