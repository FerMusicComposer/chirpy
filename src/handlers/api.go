package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type chirp struct {
	Body string `json:"body"`
}

type response struct {
	Error *string `json:"error,omitempty"`
	CleanedBody *string   `json:"cleaned_body,omitempty"`
}

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
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		w.WriteHeader(http.StatusBadRequest)
		errMsg := "something went wrong"
		res, err := json.Marshal(response{Error: &errMsg})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			return
		}
		w.Write(res)
		return
	}

	if len(newChirp.Body) > 140 {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		w.WriteHeader(http.StatusBadRequest)
		errMsg := "chirp is too long"
		res, err :=json.Marshal(response{Error: &errMsg})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			return
		}
		w.Write(res)
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

func cleanChirp(msg string) string {
	badWords:= map[string]struct{}{
		"kerfuffle": {},
		"sharbert": {},
		"fornax": {},
	}

	msgWords := strings.Split(msg, " ")

	for i, word := range msgWords {
		lowercase := strings.ToLower(word)
		if _, found := badWords[lowercase]; found {
			msgWords[i] = "****"
		}
	}

	return strings.Join(msgWords, " ")
}