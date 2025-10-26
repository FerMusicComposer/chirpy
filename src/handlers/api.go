package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/FerMusicComposer/chirpy/internal/auth"
	"github.com/FerMusicComposer/chirpy/internal/database"
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

	token, err := auth.MakeJWT(user.ID, cfg.JWTSecret, time.Hour)
	if err != nil {
		handleRequestErrors(w, "something went wrong", http.StatusInternalServerError)
		fmt.Println(fmt.Errorf("error generating JWT: %s", err))
		return
	}

	refresToken, err := auth.MakeRefreshToken()
	if err != nil {
		handleRequestErrors(w, "something went wrong", http.StatusInternalServerError)
		fmt.Println(fmt.Errorf("error generating refresh token: %s", err))
		return
	}

	err = cfg.DbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refresToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().AddDate(0, 0, 60),
	})
	if err != nil {
		handleRequestErrors(w, "something went wrong", http.StatusInternalServerError)
		fmt.Println(fmt.Errorf("error inserting refresh token on DB: %s", err))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
	res, err := json.Marshal(loginResponse{
		ID:           user.ID.String(),
		CreatedAt:    user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    user.UpdatedAt.Format(time.RFC3339),
		Email:        user.Email,
		Token:        token,
		RefreshToken: refresToken,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	w.Write(res)
}

func (cfg *ApiConfig) RefreshToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		handleRequestErrors(w, err.Error(), http.StatusBadRequest)
		fmt.Println(fmt.Errorf("error obtaining bearer: %s", err))
		return
	}

	if token == "" {
		handleRequestErrors(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	existingToken, err := cfg.DbQueries.GetRefreshToken(r.Context(), token)
	if err != nil {
		if err == sql.ErrNoRows {
			handleRequestErrors(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		handleRequestErrors(w, "something went wrong", http.StatusInternalServerError)
		fmt.Println(fmt.Errorf("error obtaining refresh token: %s", err))
		return
	}

	if existingToken.ExpiresAt.Before(time.Now()) || existingToken.RevokedAt.Valid {
		handleRequestErrors(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	newJwt, err := auth.MakeJWT(existingToken.UserID, cfg.JWTSecret, time.Hour)
	if err != nil {
		handleRequestErrors(w, "something went wrong", http.StatusInternalServerError)
		fmt.Println(fmt.Errorf("error generating JWT: %s", err))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
	res, err := json.Marshal(refreshTokenResponse{
		Token: newJwt,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	w.Write(res)
}

func (cfg *ApiConfig) RevokeToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		handleRequestErrors(w, err.Error(), http.StatusBadRequest)
		fmt.Println(fmt.Errorf("error obtaining bearer: %s", err))
		return
	}

	if token == "" {
		handleRequestErrors(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	err = cfg.DbQueries.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		handleRequestErrors(w, "something went wrong", http.StatusInternalServerError)
		fmt.Println(fmt.Errorf("error revoking refresh token: %s", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
