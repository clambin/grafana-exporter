package git

import (
	"fmt"
	"github.com/clambin/grafana-exporter/internal/writer"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"io"
	"time"
)

var _ writer.StorageHandler = &Client{}

type Client struct {
	Storage      *memory.Storage
	FS           billy.Filesystem
	Auth         transport.AuthMethod
	Repo         *git.Repository
	cloneOptions *git.CloneOptions
}

func New(url, branch, username, password string) *Client {
	auth := githttp.BasicAuth{
		Username: username,
		Password: password,
	}
	return &Client{
		Storage: memory.NewStorage(),
		FS:      memfs.New(),
		Auth:    &auth,
		cloneOptions: &git.CloneOptions{
			URL:               url,
			ReferenceName:     plumbing.NewBranchReferenceName(branch),
			Auth:              &auth,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		},
	}
}

func (c *Client) Initialize() error {
	var err error
	c.Repo, err = git.Clone(c.Storage, c.FS, c.cloneOptions)
	return err
}

func (c *Client) GetCurrent(filename string) ([]byte, error) {
	file, err := c.FS.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()
	return io.ReadAll(file)
}

func (c *Client) Mkdir(_ string) error {
	return nil
}

func (c *Client) Add(filename string, content []byte) error {
	f, err := c.FS.Create(filename)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}
	_, err = f.Write(content)
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	err = f.Close()
	if err != nil {
		return fmt.Errorf("close: %w", err)
	}
	return nil
}

func (c *Client) IsClean() (bool, error) {
	tree, err := c.Repo.Worktree()
	if err != nil {
		return false, fmt.Errorf("worktree: %w", err)
	}
	status, err := tree.Status()
	if err != nil {
		return false, fmt.Errorf("status: %w", err)
	}
	return status.IsClean(), nil
}

func (c *Client) Store() error {
	tree, err := c.Repo.Worktree()
	if err != nil {
		return fmt.Errorf("worktree: %w", err)
	}
	if _, err = tree.Add("."); err != nil {
		return fmt.Errorf("add: %w", err)
	}
	if _, err = tree.Commit("Exported Grafana dashboards", &git.CommitOptions{
		All:               false,
		AllowEmptyCommits: false,
		Author: &object.Signature{
			Name:  "grafana-exporter",
			Email: "",
			When:  time.Now(),
		},
		Committer: &object.Signature{
			Name:  "grafana-exporter",
			Email: "",
			When:  time.Now(),
		},
	}); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	if err = c.Repo.Push(&git.PushOptions{Auth: c.Auth}); err != nil {
		return fmt.Errorf("push: %w", err)
	}
	return nil
}
