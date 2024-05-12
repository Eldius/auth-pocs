package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func WithLoggingHandler() Options {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := newLoggingResponseWriter(w)
			l := slog.With(
				slog.String("request.endpoint", r.URL.String()),
				slog.String("request.host", r.Host),
				slog.String("request.remote_address", r.RemoteAddr),
				slog.String("request.body", ""), // TODO add request body
				slog.String("request.headers", fmt.Sprintf("%q", r.Header)),
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
