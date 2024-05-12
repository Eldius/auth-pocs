package middleware

import "net/http"

type Options func(handler http.Handler) http.Handler

func LoadMiddlewares(h http.Handler, opts ...Options) http.Handler {
	for _, opt := range opts {
		h = opt(h)
	}
	return h
}
