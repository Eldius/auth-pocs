package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// Start starts our server at desired port
func Start(port int) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Home)

	s := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: BasicAuthHandler(mux),
	}

	if err := s.ListenAndServe(); err != nil {
		return fmt.Errorf("starting http server: %w", err)
	}
	return nil
}

func Home(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello World!"))
}

func BasicAuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		l := slog.With(
			slog.String("endpoint", r.URL.String()),
			slog.String("host", r.Host),
			slog.String("remote_address", r.RemoteAddr),
		)

		l.DebugContext(r.Context(), "ReceivingRequest")

		u, p, ok := r.BasicAuth()
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("401 Unauthorized"))
			return
		}

		l.With(slog.String("user", u), slog.String("pass", p)).Debug("AuthData")

		next.ServeHTTP(w, r)

		slog.With(
			slog.String("response_time", time.Now().Sub(start).String()),
		).DebugContext(r.Context(), "AnsweringRequest")
	})
}
