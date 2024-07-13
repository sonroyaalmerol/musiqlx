package middlewares

import "net/http"

const (
	apiKeyHeader     = "X-Api-Key"
	apiKeyQueryParam = "apikey"
	validApiKey      = "placeholder-key"
)

func ApiKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get(apiKeyHeader)
		if apiKey == "" {
			apiKey = r.URL.Query().Get(apiKeyQueryParam)
		}

		if apiKey != validApiKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
