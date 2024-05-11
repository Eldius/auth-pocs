package api

import (
	"fmt"
	"github.com/eldius/auth-pocs/basic-auth/internal/api/middleware"
	"github.com/eldius/auth-pocs/basic-auth/internal/config"
	"github.com/eldius/auth-pocs/basic-auth/internal/persistence"
	"log/slog"
	"net/http"
)

// Start starts our server at desired port
func Start(port int) error {
	slog.Info("Starting app...")
	db := persistence.InitDB(persistence.DB(config.GetDBConfig()))

	mux := http.NewServeMux()
	mux.HandleFunc("/", Home)

	s := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: middleware.LoadMiddlewares(mux, middleware.WithLoggingHandler(), middleware.WithBasicAuthHandler(db)),
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
