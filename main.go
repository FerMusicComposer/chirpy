package main

import (
	"fmt"
	"net/http"

	"github.com/FerMusicComposer/chirpy/src/handlers"
)

func main() {
	cfg := &handlers.ApiConfig{}
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	mux.Handle("/app/", cfg.WithMetrics(http.HandlerFunc(handlers.ServeAppFiles)))
	mux.Handle("/app/assets/", cfg.WithMetrics(http.HandlerFunc(handlers.ServeAppAssets)))
	mux.HandleFunc("GET /api/healthz", handlers.GetHealthz)
	mux.HandleFunc("POST /api/validate_chirp", handlers.ValidateChirp)
	mux.HandleFunc("GET /admin/metrics", http.HandlerFunc(cfg.ServeMetrics))
	mux.HandleFunc("POST /admin/reset", http.HandlerFunc(cfg.ResetMetrics))

	fmt.Println("Listening on port 8080")
	server.ListenAndServe()

}
