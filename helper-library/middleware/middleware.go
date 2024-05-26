package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
)

type Options func(handler http.Handler) http.Handler

func LoadMiddlewares(h http.Handler, opts ...Options) http.Handler {
	slog.With(slog.String("middleware_opts", fmt.Sprintf("%+v", opts))).Error("LoadMiddlewaresStart")
	for _, opt := range opts {
		slog.With(slog.String("middleware", fmt.Sprintf("%+v", opt))).Error("LoadMiddlewaresStart")
		h = opt(h)
	}
	return h
}
