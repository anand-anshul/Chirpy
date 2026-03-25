package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/anand-anshul/chirpy/internal/auth"
	"github.com/anand-anshul/chirpy/internal/database"
	"github.com/google/uuid"
)

func handlerRediness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		ID          uuid.UUID `json:"id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Email       string    `json:"email"`
		Password    string    `json:"-"`
		IsChirpyRed bool      `json:"is_chirpy_red"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}

	user, err := cfg.dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		ID:             uuid.New(),
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprint(err))
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 501, "could not get bearer token")
		return
	}
	bearedUserID, err := auth.ValidateJWT(bearerToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "not authorised")
		return
	}

	type requestBody struct {
		Body string `json:"body"`
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
	err = decoder.Decode(&request)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	if len(request.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	cleanedBody := cleanString(request.Body)

	chirpArgs := database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: bearedUserID,
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
	chirpsSlice, err := cfg.dbQueries.GetChirps(r.Context())
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

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type requestBody struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type responseBody struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
		IsChirpyRed  bool      `json:"is_chirpy_red"`
	}

	request := requestBody{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	userStruct, err := cfg.dbQueries.GetUserByEmail(r.Context(), request.Email)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}
	match, err := auth.CheckPasswordHash(request.Password, userStruct.HashedPassword)
	if err != nil {
		respondWithError(w, 400, "could not check password hash")
		return
	}
	if !match {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	accessToken, err := auth.MakeJWT(
		userStruct.ID,
		cfg.jwtSecret,
		time.Hour,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT")
		return
	}

	refreshToken := auth.MakeRefreshToken()

	_, err = cfg.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		UserID:    userStruct.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't save refresh token")
		return
	}

	response := responseBody{
		ID:           userStruct.ID,
		CreatedAt:    userStruct.CreatedAt,
		UpdatedAt:    userStruct.UpdatedAt,
		Email:        userStruct.Email,
		Token:        accessToken,
		RefreshToken: refreshToken,
		IsChirpyRed:  userStruct.IsChirpyRed,
	}
	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't find token")
		return
	}

	user, err := cfg.dbQueries.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get user for refresh token")
		return
	}

	accessToken, err := auth.MakeJWT(
		user.ID,
		cfg.jwtSecret,
		time.Hour,
	)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate token")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't find token")
		return
	}

	_, err = cfg.dbQueries.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke session")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		ID          uuid.UUID `json:"id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Email       string    `json:"email"`
		IsChirpyRed bool      `json:"is_chirpy_red"`
	}

	request := requestBody{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithJSON(w, 401, "could not get token")
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithJSON(w, 401, "could not validate token")
		return
	}
	userStruct, err := cfg.dbQueries.UpdateUserEmail(r.Context(), database.UpdateUserEmailParams{
		ID:    userID,
		Email: request.Email,
	})
	if err != nil {
		respondWithError(w, 500, "could not update user email")
		return
	}
	hashedPass, err := auth.HashPassword(request.Password)
	if err != nil {
		respondWithError(w, 500, "could not hash password")
		return
	}
	userStruct, err = cfg.dbQueries.UpdateUserPassword(r.Context(), database.UpdateUserPasswordParams{
		ID:             userID,
		HashedPassword: hashedPass,
	})
	if err != nil {
		respondWithError(w, 500, "could not update user password")
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		ID:          userStruct.ID,
		CreatedAt:   userStruct.CreatedAt,
		UpdatedAt:   userStruct.UpdatedAt,
		Email:       userStruct.Email,
		IsChirpyRed: userStruct.IsChirpyRed,
	})
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, 400, "Invalid chirp ID")
		return
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithJSON(w, 401, "could not get token")
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithJSON(w, 401, "could not validate token")
		return
	}
	chirpStruct, err := cfg.dbQueries.GetChirp(r.Context(), chirpID)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, 404, "could not get chirp")
		return
	}
	if chirpStruct.UserID != userID {
		respondWithError(w, 403, "user not the author of chirp")
		return
	}
	err = cfg.dbQueries.DeleteChirpByUser(r.Context(), database.DeleteChirpByUserParams{
		ID:     chirpStruct.ID,
		UserID: chirpStruct.UserID,
	})
	if err != nil {
		respondWithError(w, 500, "could not delete chirp")
		return
	}
	respondWithJSON(w, 204, "chirp deleted")

}

func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type requestBody struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 401, "could not get api")
		return
	}
	if apiKey != cfg.polkaKey {
		respondWithError(w, 401, "invalid api key")
		return
	}

	request := requestBody{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&request)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	if request.Event != "user.upgraded" {
		respondWithError(w, 204, "invalid event")
		return
	}
	userID, err := uuid.Parse(request.Data.UserID)
	if err != nil {
		respondWithError(w, 500, "could not parse uuid")
		return
	}
	_, err = cfg.dbQueries.UpdateUserChirpyRed(r.Context(), userID)
	if err != nil {
		respondWithError(w, 404, "user not found")
		return
	}
	respondWithJSON(w, 204, "user upgraded successfully")

}
