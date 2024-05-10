package middleware

import "net/http"

type MiddlewareOptions func(handler http.Handler) http.Handler

func LoadMiddlewares(h http.Handler, opts ...MiddlewareOptions) http.Handler {
	for _, opt := range opts {
		h = opt(h)
	}
	return h
}
