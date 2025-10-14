package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func GetHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func ValidateChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var newChirp chirp
	err := decoder.Decode(&newChirp)

	if err != nil {
		handleRequestErrors(w, "invalid json", http.StatusBadRequest)
		return
	}

	if len(newChirp.Body) > 140 {
		handleRequestErrors(w, "chirp is too long", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
	cleanedChirp := cleanChirp(newChirp.Body)
	res, err := json.Marshal(response{CleanedBody: &cleanedChirp})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	w.Write(res)

}

func (cfg *ApiConfig) CreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req createUserRequest
	err := decoder.Decode(&req)
	if err != nil {
		handleRequestErrors(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		handleRequestErrors(w, "email is required", http.StatusBadRequest)
		return
	}

	if !validateEmail(req.Email) {
		handleRequestErrors(w, "email is invalid", http.StatusBadRequest)
		return
	}

	user, err := cfg.DbQueries.CreateUser(r.Context(), req.Email)
	if err != nil {
		handleRequestErrors(w, "error creating user", http.StatusInternalServerError)
		fmt.Println(fmt.Errorf("error creating user: %s", err))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusCreated)

	res, err := json.Marshal(createUserResponse{
		ID:        user.ID.String(),
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		Email:     user.Email,
	})
	if err != nil {
		handleRequestErrors(w, "error marshalling JSON", http.StatusInternalServerError)
		return
	}

	w.Write(res)
}
