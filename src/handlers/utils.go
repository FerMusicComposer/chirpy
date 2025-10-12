package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

func handleRequestErrors(w http.ResponseWriter, errMsg string, status int) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		w.WriteHeader(status)
		res, err := json.Marshal(response{Error: &errMsg})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			return
		}
		w.Write(res)
}

func cleanChirp(msg string) string {
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
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

func validateEmail(email string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return regex.MatchString(email)
}