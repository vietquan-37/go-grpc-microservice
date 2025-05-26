package middleware

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
)

func HealthCheckHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			payload := map[string]string{"status": "ok"}
			if err := json.NewEncoder(w).Encode(payload); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			log.Info().Msg("health check successful")
			return
		}

		handler.ServeHTTP(w, r)
	})
}
