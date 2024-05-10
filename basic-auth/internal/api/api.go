package api

import (
	"fmt"
	"github.com/eldius/auth-pocs/basic-auth/internal/repository"
	"log/slog"
	"net/http"
)

// Start starts our server at desired port
func Start(port int) error {
	slog.Info("Starting app...")
	db := repository.InitDB(repository.DB())

	mux := http.NewServeMux()
	mux.HandleFunc("/", Home)

	s := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: LoadMiddlewares(mux, WithLoggingHandler(), WithBasicAuthHandler(db)),
	}

	slog.With(slog.String("addr", s.Addr)).Info("Starting server...")
	if err := s.ListenAndServe(); err != nil {
		return fmt.Errorf("starting http server: %w", err)
	}
	return nil
}

func Home(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello World!"))
}
