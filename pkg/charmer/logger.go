package charmer

import (
	"context"
	"github.com/spf13/cobra"
	"log/slog"
)

type logCtxType string

var logCtx logCtxType = "logger"

func SetLogger(cmd *cobra.Command, logger *slog.Logger) {
	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	cmd.SetContext(context.WithValue(ctx, logCtx, logger))
}

func GetLogger(cmd *cobra.Command) *slog.Logger {
	if ctx := cmd.Context(); ctx != nil {
		if l := ctx.Value(logCtx); l != nil {
			return l.(*slog.Logger)
		}
	}
	return slog.Default()
}
