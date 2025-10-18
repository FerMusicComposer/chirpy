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

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	mux.Handle("/app/", cfg.WithMetrics(http.HandlerFunc(handlers.ServeAppFiles)))
	mux.Handle("/app/assets/", cfg.WithMetrics(http.HandlerFunc(handlers.ServeAppAssets)))
	mux.HandleFunc("GET /api/healthz", handlers.GetHealthz)
	mux.HandleFunc("POST /api/chirps", cfg.CreateChirp)
	mux.HandleFunc("POST /api/users", cfg.CreateUser)
	mux.HandleFunc("GET /admin/metrics", http.HandlerFunc(cfg.ServeMetrics))
	mux.HandleFunc("POST /admin/reset", http.HandlerFunc(cfg.ResetMetrics))

	fmt.Println("Listening on port 8080")
	server.ListenAndServe()

}
