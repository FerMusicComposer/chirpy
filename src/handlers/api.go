package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/FerMusicComposer/chirpy/internal/database"
	"github.com/google/uuid"
)

func GetHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func(cfg *ApiConfig) CreateChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	newChirp := createChirpRequest{}
	err := decoder.Decode(&newChirp)
	if err != nil {
		handleRequestErrors(w, "invalid json", http.StatusBadRequest)
		return
	}

	newChirp.Body, err = validateChirp(newChirp.Body)
	if err != nil {
		handleRequestErrors(w, err.Error(), http.StatusBadRequest)
		return
	}

	userId, err := uuid.Parse(newChirp.UserID)
	if err != nil {
		handleRequestErrors(w, "invalid user id", http.StatusBadRequest)
		return
	}

	chirp, err := cfg.DbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: newChirp.Body,
		UserID: userId,
	})
	if err != nil {
		handleRequestErrors(w, "error creating chirp", http.StatusInternalServerError)
		fmt.Println(fmt.Errorf("error creating chirp: %s", err))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusCreated)
	res, err := json.Marshal(createChirpResponse{
		ID:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt.Format(time.RFC3339),
		UpdatedAt: chirp.UpdatedAt.Format(time.RFC3339),
		Body:      chirp.Body,
		UserID:    chirp.UserID.String(),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	w.Write(res)
}

func validateChirp(body string) (string, error) {
	if len(body) > 140 {
		return "", fmt.Errorf("chirp is too long")
	}

	body = cleanChirp(body)

	return body, nil
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
