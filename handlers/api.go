package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type UserConfig struct {
	DelaySeconds        int    `json:"delay_seconds"`
	ResponseBody        string `json:"response_body"`
	StatusCode          int    `json:"status_code"`
	HangForever         bool   `json:"hang_forever"`
	HTTPMethod          string `json:"http_method"`
	RandomDelay         bool   `json:"random_delay"`
	CustomHeaders       string `json:"custom_headers"`
	FailureRate         int    `json:"failure_rate"`
	ResponseVariability string `json:"response_variability"`
	ContentType         string `json:"content_type"`
	FailureResponseBody string `json:"failure_response_body"`
	// CORS configuration
	EnableCORS       bool   `json:"enable_cors"`
	CORSAllowOrigin  string `json:"cors_allow_origin"`
	CORSAllowMethods string `json:"cors_allow_methods"`
	CORSAllowHeaders string `json:"cors_allow_headers"`
	CORSMaxAge       int    `json:"cors_max_age"`
	CORSAllowCreds   bool   `json:"cors_allow_credentials"`
}

func HandleAPIRequest(w http.ResponseWriter, r *http.Request, cfg UserConfig) {
	if cfg.HTTPMethod != "" && r.Method != cfg.HTTPMethod {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Apply CORS headers if enabled
	if cfg.EnableCORS {
		// Set Access-Control-Allow-Origin header
		if cfg.CORSAllowOrigin != "" {
			w.Header().Set("Access-Control-Allow-Origin", cfg.CORSAllowOrigin)
		} else {
			// Default to all origins if not specified
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			// Set Access-Control-Allow-Methods header
			if cfg.CORSAllowMethods != "" {
				w.Header().Set("Access-Control-Allow-Methods", cfg.CORSAllowMethods)
			} else {
				// Default methods if not specified
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			}

			// Set Access-Control-Allow-Headers header
			if cfg.CORSAllowHeaders != "" {
				w.Header().Set("Access-Control-Allow-Headers", cfg.CORSAllowHeaders)
			} else {
				// Default headers if not specified
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			}

			// Set Access-Control-Max-Age header if specified
			if cfg.CORSMaxAge > 0 {
				w.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", cfg.CORSMaxAge))
			}

			// Set Access-Control-Allow-Credentials header if enabled
			if cfg.CORSAllowCreds {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			// Respond to preflight request with 204 No Content
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// For non-preflight requests, still set credentials header if enabled
		if cfg.CORSAllowCreds {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
	}

	if cfg.ContentType != "" {
		w.Header().Set("Content-Type", cfg.ContentType)
	}

	if cfg.CustomHeaders != "" {
		headers := strings.Split(cfg.CustomHeaders, "\n")
		for _, header := range headers {
			parts := strings.SplitN(header, ":", 2)
			if len(parts) == 2 {
				w.Header().Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
			}
		}
	}

	if cfg.FailureRate > 0 && rand.Intn(100) < cfg.FailureRate {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, cfg.FailureResponseBody)
		return
	}

	if cfg.HangForever {
		select {}
	}

	delay := time.Duration(cfg.DelaySeconds) * time.Second
	if cfg.RandomDelay {
		variation := rand.Float64() - 0.5
		delay = time.Duration(float64(delay) * (1 + variation))
	}

	time.Sleep(delay)
	w.WriteHeader(cfg.StatusCode)
	fmt.Fprint(w, cfg.ResponseBody)
}
