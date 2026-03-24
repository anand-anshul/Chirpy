package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(response)

	return nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) error {
	return respondWithJSON(w, code, map[string]string{"error": msg})
}

func cleanString(s string) string {
	sSlice := strings.Split(s, " ")
	for i, word := range sSlice {
		if strings.ToLower(word) == "kerfuffle" {
			sSlice[i] = "****"
		} else if strings.ToLower(word) == "fornax" {
			sSlice[i] = "****"
		} else if strings.ToLower(word) == "sharbert" {
			sSlice[i] = "****"
		}
	}
	cleanString := strings.Join(sSlice, " ")
	return cleanString
}
