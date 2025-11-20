package handler

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"ool/api/schema"

	"golang.org/x/time/rate"
)

var clientLimiters = make(map[string]*rate.Limiter)
var limiterMutex sync.Mutex

func getLimiter(apiKey string) *rate.Limiter {
	limiterMutex.Lock()
	defer limiterMutex.Unlock()

	limiter, exists := clientLimiters[apiKey]
	if !exists {
		limiter = rate.NewLimiter(rate.Every(time.Second/10), 20)
		clientLimiters[apiKey] = limiter
	}
	return limiter
}

// respondJSON writes a JSON response with the specified status code and payload.
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Internal Server Error: JSON marshalling failed", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

// respondError writes a standard error JSON response with the specified code and message,
func respondError(w http.ResponseWriter, code int, message string) {
	errResponse := schema.ApiResponse{
		Code:    code,
		Type:    "error",
		Message: message,
	}
	respondJSON(w, code, errResponse)
}

// APIKeyMiddleware checks for the required 'api_key' header
// Hardcode value: apitest
func APIKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const expectedAPIKey = "apitest"
		apiKey := r.Header.Get("api_key")

		if apiKey != expectedAPIKey {
			respondError(w, http.StatusUnauthorized, "Invalid API Key")
			return
		}

		limiter := getLimiter(apiKey)

		if !limiter.Allow() {
			w.Header().Set("Retry-After", "1")
			respondError(w, http.StatusTooManyRequests, "Too many request")
			return
		}

		next.ServeHTTP(w, r)
	})
}
