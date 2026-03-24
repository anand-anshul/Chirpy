package main

import (
	"encoding/json"
	"net/http"
)

type requestBody struct {
	Body string `json:"body"`
}
type responseBody struct {
	CleanedBody string `json:"cleaned_body"`
}

func handlerRediness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	params := requestBody{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	cleanedBody := cleanString(params.Body)

	response := responseBody{
		CleanedBody: cleanedBody,
	}
	respondWithJSON(w, 200, response)
}
