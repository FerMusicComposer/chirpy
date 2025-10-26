package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/FerMusicComposer/chirpy/internal/auth"
	"github.com/FerMusicComposer/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) UpdateUserSubscriptionWebhook(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		handleRequestErrors(w, err.Error(), http.StatusUnauthorized)
		fmt.Println(fmt.Errorf("error obtaining api key: %s", err))
		return
	}

	if apiKey == "" || apiKey != cfg.PolkaKey {
		handleRequestErrors(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var req webhookRequest
	err = decoder.Decode(&req)
	if err != nil {
		handleRequestErrors(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = cfg.DbQueries.UpdateUserSubscription(r.Context(), database.UpdateUserSubscriptionParams{
		ID:          uuid.MustParse(req.Data.UserId),
		IsChirpyRed: true,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			handleRequestErrors(w, "user not found", http.StatusNotFound)
			return
		}

		handleRequestErrors(w, "something went wrong", http.StatusInternalServerError)
		fmt.Println(fmt.Errorf("error updating user: %s", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
