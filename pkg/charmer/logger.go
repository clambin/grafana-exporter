package charmer

import (
	"context"
	"github.com/spf13/cobra"
	"log/slog"
)

type logCtxType string

var logCtx logCtxType = "logger"

func SetLogger(cmd *cobra.Command, logger *slog.Logger) {
	ctx := context.WithValue(cmd.Context(), logCtx, logger)
	cmd.SetContext(ctx)
}

func GetLogger(cmd *cobra.Command) *slog.Logger {
	if l := cmd.Context().Value(logCtx); l != nil {
		return l.(*slog.Logger)
	}
	return slog.Default()
}
