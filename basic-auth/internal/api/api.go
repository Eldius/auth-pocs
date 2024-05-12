package api

import (
	"fmt"
	"github.com/eldius/auth-pocs/basic-auth/internal/api/middleware"
	"github.com/eldius/auth-pocs/basic-auth/internal/auth"
	"github.com/eldius/auth-pocs/basic-auth/internal/persistence"
	helpermiddleware "github.com/eldius/auth-pocs/helper-library/middleware"
	"log/slog"
	"net/http"
)

// Start starts our server at desired port
func Start(port int) error {
	slog.Info("Starting app...")
	db := persistence.InitDB()

	mux := http.NewServeMux()
	mux.HandleFunc("/", Home)

	s := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: helpermiddleware.LoadMiddlewares(mux, helpermiddleware.WithLoggingHandler(), middleware.WithBasicAuthHandler(db)),
	}
	slog.With(slog.String("addr", s.Addr)).Info("Starting server...")
	if err := s.ListenAndServe(); err != nil {
		return fmt.Errorf("starting http server: %w", err)
	}
	return nil
}

func Home(w http.ResponseWriter, r *http.Request) {
	u := auth.UserFromContext(r.Context())
	_, _ = w.Write([]byte(fmt.Sprintf("Hello, %s!", u.User)))
}
