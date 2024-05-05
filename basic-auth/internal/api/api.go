package api

import (
	"fmt"
	"net/http"
)

// Start starts our server at desired port
func Start(port int) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Home)

	s := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	if err := s.ListenAndServe(); err != nil {
		return fmt.Errorf("starting http server: %w", err)
	}
	return nil
}

func Home(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello World!"))
}
