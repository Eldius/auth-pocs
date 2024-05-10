package middleware

import (
	"fmt"
	"github.com/eldius/auth-pocs/basic-auth/internal/auth"
	"github.com/eldius/auth-pocs/basic-auth/internal/repository"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"net/http"
)

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
