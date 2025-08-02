package middleware

import (
	"net/http"
	"os"
	"strings"
)

// CORSMiddleware adds CORS headers to responses based on environment variables
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if CORS is enabled
		corsEnabled := os.Getenv("CORS_ENABLED")
		if corsEnabled != "true" {
			// CORS is not enabled, just call the next handler
			next.ServeHTTP(w, r)
			return
		}

		// Get CORS configuration from environment variables
		allowOrigin := os.Getenv("CORS_ALLOW_ORIGIN")
		if allowOrigin == "" {
			allowOrigin = "*" // Default to all origins if not specified
		}

		allowMethods := os.Getenv("CORS_ALLOW_METHODS")
		if allowMethods == "" {
			allowMethods = "GET, POST, PUT, DELETE, OPTIONS" // Default methods
		}

		allowHeaders := os.Getenv("CORS_ALLOW_HEADERS")
		if allowHeaders == "" {
			allowHeaders = "Content-Type, Authorization" // Default headers
		}

		maxAge := os.Getenv("CORS_MAX_AGE")
		if maxAge == "" {
			maxAge = "86400" // Default to 24 hours
		}

		allowCredentials := os.Getenv("CORS_ALLOW_CREDENTIALS")
		exposeHeaders := os.Getenv("CORS_EXPOSE_HEADERS")

		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		w.Header().Set("Access-Control-Allow-Methods", allowMethods)
		w.Header().Set("Access-Control-Allow-Headers", allowHeaders)
		w.Header().Set("Access-Control-Max-Age", maxAge)

		if allowCredentials == "true" {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if exposeHeaders != "" {
			w.Header().Set("Access-Control-Expose-Headers", exposeHeaders)
		}

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}