package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"net/http"
	"time"
)

var (
	AuthUserNotFoundErr = errors.New("user lookup")
)

type ctxAuthUserInfo struct {
}

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
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u, p, ok := r.BasicAuth()
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte("401 Unauthorized"))
				return
			}
			rows, err := db.NamedQueryContext(r.Context(), `select * from auth_users where user = :user`, map[string]interface{}{
				"user": u,
			})
			if err != nil {
				err = fmt.Errorf("user lookup: %w", err)
				slog.With("error", err).Error("AuthData")
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte("401 Unauthorized"))
				return
			}
			defer func() {
				_ = rows.Close()
			}()

			if !rows.Next() {
				slog.With("error", AuthUserNotFoundErr).Error("AuthData")
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte("401 Unauthorized"))
				return
			}

			usrInfoMap := map[string]interface{}{}

			if err := rows.MapScan(usrInfoMap); err != nil {
				err = fmt.Errorf("rows mapping: %w", err)
				slog.With("error", err).Error("AuthData")
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte("401 Unauthorized"))
				return
			}

			if pass, ok := usrInfoMap["pass"]; !ok || pass != p {
				slog.Error("AuthDataNotAuthorized")
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte("401 Unauthorized"))
				return
			}

			r.WithContext(context.WithValue(r.Context(), ctxAuthUserInfo{}, usrInfoMap["user"]))

			slog.With(slog.String("user", u), slog.String("pass", p)).Debug("AuthData")
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
