package handlers

import (
	"fmt"
	"net/http"
)

func (cfg *ApiConfig) WithMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *ApiConfig) ServeMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
	res := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.FileServerHits.Load())
	w.Write([]byte(res))
}

func (cfg *ApiConfig) ResetMetrics(w http.ResponseWriter, r *http.Request) {
	if cfg.Environment == "dev" {
		err := cfg.DbQueries.DeleteAllUsers(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	cfg.FileServerHits.Store(0)
}
