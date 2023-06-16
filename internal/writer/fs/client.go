package fs

import (
	"github.com/clambin/grafana-exporter/internal/writer"
	"golang.org/x/exp/slog"
	"os"
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

func (c *Client) Mkdir(s string) error {
	slog.Debug("creating local directory")
	return os.MkdirAll(s, 0755)
}

func (c *Client) Add(s string, bytes []byte) error {
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
