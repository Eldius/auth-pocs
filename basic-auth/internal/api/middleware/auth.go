package middleware

import (
	"fmt"
	"github.com/eldius/auth-pocs/basic-auth/internal/auth"
	"github.com/eldius/auth-pocs/basic-auth/internal/persistence/repository"
	"github.com/eldius/auth-pocs/helper-library/middleware"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"net/http"
	"reflect"
	"strings"
)

func WithBasicAuthHandler(db *sqlx.DB, publicEndpoints ...string) middleware.Options {
	repo := repository.NewUserRepository(db)
	svc := auth.NewAuthService(repo)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("%+v\n", w)
			fmt.Printf("type of w:%+v\n", reflect.TypeOf(w).String())
			slog.With(
				slog.String("url", r.URL.String()),
				slog.String("publicEndpoints", strings.Join(publicEndpoints, ",")),
			).Info("AuthDataStart")
			fmt.Println("AuthDataStart")

			for _, endpoint := range publicEndpoints {
				if endpoint == r.URL.String() {
					slog.With(
						slog.String("url", r.URL.String()),
						slog.String("publicEndpoints", strings.Join(publicEndpoints, ",")),
					).Info("AuthDataSkip")
					next.ServeHTTP(w, r)
					return
				}
			}
			//if slices.Contains(publicEndpoints, r.URL.Path) {
			//	slog.With(
			//		slog.String("url", r.URL.String()),
			//		slog.String("publicEndpoints", strings.Join(publicEndpoints, ",")),
			//	).Info("AuthDataSkip")
			//	next.ServeHTTP(w, r)
			//	return
			//}
			u, p, ok := r.BasicAuth()
			if !ok {
				unauthorized(w)
				return
			}
			_, ctx, err := svc.AuthenticateUser(r.Context(), u, p)
			if err != nil {
				unauthorized(w)
				err = fmt.Errorf("invalid user: %w", err)
				slog.With("error", err).Error("AuthData")
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func unauthorized(w http.ResponseWriter) {
	w.Header().Add("WWW-Authenticate", "Basic realm=\"app access\"")

	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte("401 Unauthorized"))
}
