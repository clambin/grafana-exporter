package fs

import (
	"errors"
	"fmt"
	"github.com/clambin/grafana-exporter/internal/writer"
	"os"
	"path"
)

var _ writer.StorageHandler = &Client{}

type file struct {
	filename string
	content  []byte
}

type Client struct {
	files []file
}

func (c *Client) Initialize() error {
	c.files = nil
	return nil
}

func (c *Client) GetCurrent(s string) ([]byte, error) {
	return os.ReadFile(s)
}

func (c *Client) Add(s string, bytes []byte) error {
	c.files = append(c.files, file{filename: s, content: bytes})
	return nil
}

func (c *Client) IsClean() (bool, error) {
	return len(c.files) == 0, nil
}

func (c *Client) Store(_ string) error {
	for _, f := range c.files {
		dirname := path.Dir(f.filename)
		if err := os.MkdirAll(dirname, 0755); err != nil && !errors.Is(err, os.ErrExist) {
			return fmt.Errorf("mkdir %s: %w", dirname, err)
		}
		if err := os.WriteFile(f.filename, f.content, 0644); err != nil {
			return err
		}
	}
	return c.Initialize()
}
