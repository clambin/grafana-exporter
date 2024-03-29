package charmer_test

import (
	"context"
	"github.com/clambin/grafana-exporter/pkg/charmer"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	var cmd cobra.Command

	logger := charmer.GetLogger(&cmd)
	if logger != slog.Default() {
		t.Errorf("logger should return default logger")
	}

	cmd.SetContext(context.WithValue(context.Background(), "logger", nil)) //nolint:all
	if logger != slog.Default() {
		t.Errorf("logger should return default logger")
	}

	l := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{}))
	charmer.SetLogger(&cmd, l)

	logger = charmer.GetLogger(&cmd)
	if logger != l {
		t.Errorf("logger should return new logger")
	}
}
