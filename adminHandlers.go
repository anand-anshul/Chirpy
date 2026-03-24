package main

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/anand-anshul/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	htmlString := `
		<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>
	`
	reqCountHtml := fmt.Sprintf(htmlString, cfg.fileserverHits.Load())
	w.Write([]byte(reqCountHtml))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, 403, "Forbidden")
		return
	}
	cfg.fileserverHits.Store(0)

	err := cfg.dbQueries.DeleteAllChirps(r.Context())
	if err != nil {
		respondWithError(w, 500, "could not delete all chirps")
		return
	}

	err = cfg.dbQueries.DeleteAllUsers(r.Context())
	if err != nil {
		respondWithError(w, 500, "could not delete all users")
		return
	}
	w.WriteHeader(http.StatusOK)

}
