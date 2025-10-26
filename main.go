package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/FerMusicComposer/chirpy/internal/database"
	"github.com/FerMusicComposer/chirpy/src/handlers"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dbQueries := database.New(db)
	defer db.Close()

	cfg := &handlers.ApiConfig{}
	cfg.DbQueries = dbQueries
	cfg.Environment = os.Getenv("PLATFORM")
	cfg.JWTSecret = os.Getenv("JWT_SECRET")

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// General
	mux.Handle("/app/", cfg.WithMetrics(http.HandlerFunc(handlers.ServeAppFiles)))
	mux.Handle("/app/assets/", cfg.WithMetrics(http.HandlerFunc(handlers.ServeAppAssets)))
	mux.HandleFunc("GET /api/healthz", handlers.GetHealthz)

	// Auth
	mux.HandleFunc("POST /api/login", cfg.Login)
	mux.HandleFunc("POST /api/refresh", cfg.RefreshToken)
	mux.HandleFunc("POST /api/revoke", cfg.RevokeToken)

	// Chirps
	mux.HandleFunc("POST /api/chirps", cfg.CreateChirp)
	mux.HandleFunc("GET /api/chirps", cfg.GetChirps)
	mux.HandleFunc("GET /api/chirps/{id}", cfg.GetChirp)
	mux.HandleFunc("DELETE /api/chirps/{id}", cfg.DeleteChirp)

	// Users
	mux.HandleFunc("POST /api/users", cfg.CreateUser)
	mux.HandleFunc("PUT /api/users", cfg.UpdateUser)

	// Admin
	mux.HandleFunc("GET /admin/metrics", http.HandlerFunc(cfg.ServeMetrics))
	mux.HandleFunc("POST /admin/reset", http.HandlerFunc(cfg.ResetMetrics))

	fmt.Println("Listening on port 8080")
	server.ListenAndServe()

}
