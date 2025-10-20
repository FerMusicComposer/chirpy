package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/FerMusicComposer/chirpy/internal/auth"
)

func GetHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *ApiConfig) Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	req := loginRequest{}
	err := decoder.Decode(&req)
	if err != nil {
		handleRequestErrors(w, "invalid json", http.StatusBadRequest)
		return
	}
	
	if req.Email == "" {
		handleRequestErrors(w, "email is required", http.StatusBadRequest)
		return
	}

	user, err := cfg.DbQueries.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			handleRequestErrors(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		handleRequestErrors(w, "something went wrong", http.StatusInternalServerError)
		fmt.Println(fmt.Errorf("error getting user: %s", err))
		return
	}

	matched, err := auth.CheckPasswordHash(req.Password, user.HashedPassword)
	if err != nil {
		handleRequestErrors(w, "something went wrong", http.StatusInternalServerError)
		fmt.Println(fmt.Errorf("error checking password: %s", err))
		return
	}

	if !matched {
		handleRequestErrors(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var expirationTime time.Duration
	if req.ExpiresInSeconds == 0 || req.ExpiresInSeconds > 3600 {
		expirationTime = time.Hour
	} else {
		expirationTime = time.Second * time.Duration(req.ExpiresInSeconds)
	}

	token, err := auth.MakeJWT(user.ID, cfg.JWTSecret, expirationTime)
	if err != nil {
		handleRequestErrors(w, "something went wrong", http.StatusInternalServerError)
		fmt.Println(fmt.Errorf("error making JWT: %s", err))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
	res, err := json.Marshal(loginResponse{
		ID: user.ID.String(),
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		Email: user.Email,
		Token: token,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	w.Write(res)
}