package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/anand-anshul/chirpy/internal/database"
	"github.com/google/uuid"
)

func handlerRediness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
// 	defer r.Body.Close()

// 	type requestBody struct {
// 		Body string `json:"body"`
// 	}
// 	type responseBody struct {
// 		CleanedBody string `json:"cleaned_body"`
// 	}
// 	request := requestBody{}

// 	decoder := json.NewDecoder(r.Body)
// 	err := decoder.Decode(&request)
// 	if err != nil {
// 		respondWithError(w, 500, "Something went wrong")
// 		return
// 	}

// 	if len(request.Body) > 140 {
// 		respondWithError(w, 400, "Chirp is too long")
// 		return
// 	}

// 	cleanedBody := cleanString(request.Body)

// 	response := responseBody{
// 		CleanedBody: cleanedBody,
// 	}
// 	respondWithJSON(w, 200, response)
// }

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type requestBody struct {
		Email string `json:"email"`
	}
	type responseBody struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}
	request := requestBody{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		respondWithError(w, 500, "could not decode request")
		return
	}
	userStruct, err := cfg.dbQueries.CreateUser(r.Context(), request.Email)
	if err != nil {
		respondWithError(w, 500, "could not create user")
		return
	}

	response := responseBody{
		ID:        userStruct.ID,
		CreatedAt: userStruct.CreatedAt,
		UpdatedAt: userStruct.UpdatedAt,
		Email:     userStruct.Email,
	}

	respondWithJSON(w, http.StatusCreated, response)

}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type requestBody struct {
		Body   string `json:"body"`
		UserID string `json:"user_id"`
	}
	type responseBody struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}
	request := requestBody{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	if len(request.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	cleanedBody := cleanString(request.Body)

	userID, err := uuid.Parse(request.UserID)
	if err != nil {
		respondWithError(w, 400, "Invalid user ID")
		return
	}

	chirpArgs := database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: userID,
	}

	chirpStruct, err := cfg.dbQueries.CreateChirp(r.Context(), chirpArgs)
	if err != nil {
		respondWithError(w, 500, "could not create user")
		return
	}

	response := responseBody{
		ID:        chirpStruct.ID,
		CreatedAt: chirpStruct.CreatedAt,
		UpdatedAt: chirpStruct.UpdatedAt,
		Body:      chirpStruct.Body,
		UserID:    chirpArgs.UserID,
	}

	respondWithJSON(w, http.StatusCreated, response)
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type responseBody []struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}
	chirpsSlice, err := cfg.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, 500, "could not get chirps")
		return
	}
	response := responseBody{}

	for _, chirp := range chirpsSlice {
		response = append(response, struct {
			ID        uuid.UUID `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Body      string    `json:"body"`
			UserID    uuid.UUID `json:"user_id"`
		}{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type responseBody struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, 400, "Invalid chirp ID")
		return
	}
	chirpStruct, err := cfg.dbQueries.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, 404, "could not get chirps")
		return
	}

	response := responseBody{
		ID:        chirpStruct.ID,
		CreatedAt: chirpStruct.CreatedAt,
		UpdatedAt: chirpStruct.UpdatedAt,
		Body:      chirpStruct.Body,
		UserID:    chirpStruct.UserID,
	}

	respondWithJSON(w, http.StatusOK, response)

}
