package handlers

import "net/http"

func ServeAppFiles(w http.ResponseWriter, r *http.Request) {
	fileServer := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	fileServer.ServeHTTP(w, r)
}

func ServeAppAssets(w http.ResponseWriter, r *http.Request) {
	fileServer := http.StripPrefix("/app/assets/", http.FileServer(http.Dir("./src/assets")))
	fileServer.ServeHTTP(w, r)
}
