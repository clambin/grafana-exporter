package cli

import (
	"errors"
	"fmt"
	"github.com/clambin/grafana-exporter/internal/writer"
	"github.com/clambin/grafana-exporter/internal/writer/fs"
	"github.com/clambin/grafana-exporter/internal/writer/git"
	"github.com/spf13/viper"
)

func makeWriter() (*writer.Writer, error) {
	var storage writer.StorageHandler
	mode := viper.GetString("mode")
	switch mode {
	case "local":
		storage = &fs.Client{}
	case "git":
		url := viper.GetString("git.url")
		if url == "" {
			return nil, errors.New("missing git url")
		}
		branch := viper.GetString("git.branch")
		username := viper.GetString("git.username")
		if username == "" {
			return nil, errors.New("missing git username")
		}
		password := viper.GetString("git.password")
		if password == "" {
			return nil, errors.New("missing git password/token")
		}
		storage = git.New(url, branch, username, password)
	default:
		return nil, fmt.Errorf("invalid output mode: '%s'", mode)
	}
	return &writer.Writer{StorageHandler: storage, BaseDirectory: viper.GetString("out")}, nil
}
