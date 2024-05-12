package logging

import (
	"context"
	"log/slog"
)

func ContextLogger(ctx context.Context) *slog.Logger {
	return slog.With()
}
