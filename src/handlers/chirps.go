package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/FerMusicComposer/chirpy/internal/database"
	"github.com/google/uuid"
)

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

func(cfg *ApiConfig) GetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.DbQueries.GetAllChirps(r.Context())
	if err != nil {
		handleRequestErrors(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(fmt.Errorf("error getting chirps: %s", err))
		return
	}

	resp := make([]createChirpResponse, len(chirps))
	for i, chirp := range chirps {
		resp[i] = createChirpResponse{
			ID:        chirp.ID.String(),
			CreatedAt: chirp.CreatedAt.Format(time.RFC3339),
			UpdatedAt: chirp.UpdatedAt.Format(time.RFC3339),
			Body:      chirp.Body,
			UserID:    chirp.UserID.String(),
		}
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