package fs

import (
	"errors"
	"fmt"
	"github.com/clambin/grafana-exporter/internal/writer"
	"os"
	"path"
)

var _ writer.StorageHandler = &Client{}

type fileToStore struct {
	filename string
	content  []byte
}

type Client struct {
	files []fileToStore
}

func (c *Client) Initialize() error {
	c.files = []fileToStore{}
	return nil
}

func (c *Client) GetCurrent(s string) ([]byte, error) {
	return os.ReadFile(s)
}

func (c *Client) Add(s string, bytes []byte) error {
	dirname := path.Dir(s)
	if err := os.MkdirAll(dirname, 0755); err != nil && !errors.Is(err, os.ErrExist) {
		return fmt.Errorf("mkdir %s: %w", dirname, err)
	}

	c.files = append(c.files, fileToStore{filename: s, content: bytes})
	return nil
}

func (c *Client) IsClean() (bool, error) {
	return len(c.files) == 0, nil
}

func (c *Client) Store() error {
	for _, file := range c.files {
		if err := os.WriteFile(file.filename, file.content, 0644); err != nil {
			return err
		}
	}
	return c.Initialize()
}
