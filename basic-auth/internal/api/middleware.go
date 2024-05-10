package api

import (
	"fmt"
	"github.com/eldius/auth-pocs/basic-auth/internal/auth"
	"github.com/eldius/auth-pocs/basic-auth/internal/repository"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func WithLoggingHandler() MiddlewareOptions {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := NewLoggingResponseWriter(w)
			l := slog.With(
				slog.String("request.endpoint", r.URL.String()),
				slog.String("request.host", r.Host),
				slog.String("request.remote_address", r.RemoteAddr),
				slog.String("request.body", ""), // TODO add request body
			)

			l.DebugContext(r.Context(), "ReceivingRequest")

			next.ServeHTTP(ww, r)

			l.With(
				slog.String("response.response_time", time.Since(start).String()),
				slog.Int("response.status_code", ww.statusCode),
				slog.String("response.body", ""), // TODO add response body
			).DebugContext(r.Context(), "AnsweringRequest")
		})
	}
}

func WithBasicAuthHandler(db *sqlx.DB) MiddlewareOptions {
	repo := repository.NewUserRepository(db)
	svc := auth.NewAuthService(repo)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u, p, ok := r.BasicAuth()
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte("401 Unauthorized"))
				return
			}
			_, ctx, err := svc.AuthenticateUser(r.Context(), u, p)
			if err != nil {
				err = fmt.Errorf("invalid user: %w", err)
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte("401 Unauthorized"))
				slog.With("error", err).Error("AuthData")
				return
			}

			r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

type MiddlewareOptions func(handler http.Handler) http.Handler

func LoadMiddlewares(h http.Handler, opts ...MiddlewareOptions) http.Handler {
	for _, opt := range opts {
		h = opt(h)
	}
	return h
}
