package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/FerMusicComposer/chirpy/internal/auth"
	"github.com/FerMusicComposer/chirpy/internal/database"
)

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

	if req.Password == "" {
		handleRequestErrors(w, "password is required", http.StatusBadRequest)
		return
	}

	hashedPwd, err := auth.HashPassword(req.Password)
	if err != nil {
		handleRequestErrors(w, "something went wrong", http.StatusInternalServerError)
		fmt.Println(fmt.Errorf("error hashing password: %s", err))
		return
	}

	user, err := cfg.DbQueries.CreateUser(r.Context(), database.CreateUserParams{
		Email: req.Email,
		HashedPassword: hashedPwd,
	})
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
