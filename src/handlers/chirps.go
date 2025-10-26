package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/FerMusicComposer/chirpy/internal/auth"
	"github.com/FerMusicComposer/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) CreateChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	newChirp := createChirpRequest{}
	err := decoder.Decode(&newChirp)
	if err != nil {
		handleRequestErrors(w, "invalid json", http.StatusBadRequest)
		fmt.Println(fmt.Errorf("error decoding json: %s", err))
		return
	}

	newChirp.Body, err = validateChirp(newChirp.Body)
	if err != nil {
		handleRequestErrors(w, err.Error(), http.StatusBadRequest)
		fmt.Println(fmt.Errorf("error validating chirp: %s", err))
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		handleRequestErrors(w, err.Error(), http.StatusBadRequest)
		fmt.Println(fmt.Errorf("error obtaining bearer: %s", err))
		return
	}

	jwtUserId, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		handleRequestErrors(w, "unauthorized", http.StatusUnauthorized)
		fmt.Println(fmt.Errorf("error validating jwt: %s", err))
		return
	}

	chirp, err := cfg.DbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   newChirp.Body,
		UserID: jwtUserId,
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
		baseModel: baseModel{
			ID:        chirp.ID.String(),
			CreatedAt: chirp.CreatedAt.Format(time.RFC3339),
			UpdatedAt: chirp.UpdatedAt.Format(time.RFC3339),
		},
		chirpData: chirpData{
			Body:   chirp.Body,
			UserID: chirp.UserID.String(),
		},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	w.Write(res)
}

func (cfg *ApiConfig) GetChirps(w http.ResponseWriter, r *http.Request) {
	var chirps []database.Chirp
	var err error

	userId := r.URL.Query().Get("author_id")

	if userId != "" {
		chirps, err = cfg.DbQueries.GetChirpsByAuthor(r.Context(), uuid.MustParse(userId))
	} else {
		chirps, err = cfg.DbQueries.GetAllChirps(r.Context())
	}

	if err != nil {
		handleRequestErrors(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(fmt.Errorf("error getting chirps: %s", err))
		return
	}

	resp := make([]createChirpResponse, len(chirps))
	for i, chirp := range chirps {
		resp[i] = createChirpResponse{
			baseModel: baseModel{
				ID:        chirp.ID.String(),
				CreatedAt: chirp.CreatedAt.Format(time.RFC3339),
				UpdatedAt: chirp.UpdatedAt.Format(time.RFC3339),
			},
			chirpData: chirpData{
				Body:   chirp.Body,
				UserID: chirp.UserID.String(),
			}}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
	res, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	w.Write(res)
}

func (cfg *ApiConfig) GetChirp(w http.ResponseWriter, r *http.Request) {
	chirpId := r.PathValue("id")

	chirp, err := cfg.DbQueries.GetChirp(r.Context(), uuid.MustParse(chirpId))
	if err != nil {
		if err == sql.ErrNoRows {
			handleRequestErrors(w, "chirp not found", http.StatusNotFound)
			return
		}

		handleRequestErrors(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(fmt.Errorf("error getting chirp: %s", err))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
	res, err := json.Marshal(createChirpResponse{
		baseModel: baseModel{
			ID:        chirp.ID.String(),
			CreatedAt: chirp.CreatedAt.Format(time.RFC3339),
			UpdatedAt: chirp.UpdatedAt.Format(time.RFC3339),
		},
		chirpData: chirpData{
			Body:   chirp.Body,
			UserID: chirp.UserID.String(),
		},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	w.Write(res)
}

func (cfg *ApiConfig) DeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpId := r.PathValue("id")

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		handleRequestErrors(w, err.Error(), http.StatusUnauthorized)
		fmt.Println(fmt.Errorf("error obtaining bearer: %s", err))
		return
	}

	jwtUserId, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		handleRequestErrors(w, "unauthorized", http.StatusUnauthorized)
		fmt.Println(fmt.Errorf("error validating jwt: %s", err))
		return
	}

	chirp, err := cfg.DbQueries.GetChirp(r.Context(), uuid.MustParse(chirpId))
	if err != nil {
		if err == sql.ErrNoRows {
			handleRequestErrors(w, "chirp not found", http.StatusNotFound)
			return
		}

		handleRequestErrors(w, "something went wrong", http.StatusInternalServerError)
		fmt.Println(fmt.Errorf("error getting chirp: %s", err))
		return
	}

	if chirp.UserID != jwtUserId {
		handleRequestErrors(w, "forbidden", http.StatusForbidden)
		return
	}

	err = cfg.DbQueries.DeleteChirp(r.Context(), uuid.MustParse(chirpId))
	if err != nil {
		handleRequestErrors(w, "error deleting chirp", http.StatusInternalServerError)
		fmt.Println(fmt.Errorf("error deleting chirp: %s", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
