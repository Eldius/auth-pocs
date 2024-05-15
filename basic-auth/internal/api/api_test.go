package api

import (
	"github.com/eldius/auth-pocs/basic-auth/internal/api/middleware"
	"github.com/eldius/auth-pocs/basic-auth/internal/persistence"
	"github.com/eldius/auth-pocs/helper-library/logging"
	hmiddleware "github.com/eldius/auth-pocs/helper-library/middleware"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHome(t *testing.T) {

	assert.Nil(t, logging.SetupLogs("basic-auth-test", true))
	db := persistence.InitDB()
	defer func() {
		_ = db.Close()
	}()

	t.Run("given an unauthenticated request should return 401", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/", Home)
		server := httptest.NewServer(hmiddleware.LoadMiddlewares(mux, middleware.WithBasicAuthHandler(db)))
		defer server.Close()

		resp, err := http.Get(server.URL + "/")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("given an authenticated request should return 200", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/", Home)
		server := httptest.NewServer(hmiddleware.LoadMiddlewares(mux, hmiddleware.WithLoggingHandler(), middleware.WithBasicAuthHandler(db)))
		defer server.Close()

		req, err := http.NewRequest(http.MethodGet, server.URL+"/", nil)
		req.SetBasicAuth("root", "12345")

		t.Logf("auth_header: %+v", req.Header)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("given request using wrong password should return 401", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/", Home)
		server := httptest.NewServer(hmiddleware.LoadMiddlewares(mux, hmiddleware.WithLoggingHandler(), middleware.WithBasicAuthHandler(db)))
		defer server.Close()

		req, err := http.NewRequest(http.MethodGet, server.URL+"/", nil)
		req.SetBasicAuth("root", "123")

		t.Logf("auth_header: %+v", req.Header)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("given request using wrong user should return 401", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/", Home)
		server := httptest.NewServer(hmiddleware.LoadMiddlewares(mux, hmiddleware.WithLoggingHandler(), middleware.WithBasicAuthHandler(db)))
		defer server.Close()

		req, err := http.NewRequest(http.MethodGet, server.URL+"/", nil)
		req.SetBasicAuth("root123", "12345")

		t.Logf("auth_header: %+v", req.Header)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
