package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type UserConfig struct {
	HTTPMethod string `json:"http_method"`

	DelaySeconds        int    `json:"delay_seconds"`
	HangForever         bool   `json:"hang_forever"`
	FailureRate         int    `json:"failure_rate"`
	CustomHeaders       string `json:"custom_headers"`
	FailureResponseBody string `json:"failure_response_body"`

	ResponseBody string `json:"response_body"`
	StatusCode   int    `json:"status_code"`
	ContentType  string `json:"content_type"`
}

func HandleAPIRequest(w http.ResponseWriter, r *http.Request, config UserConfig) {
	if config.HTTPMethod != "" && r.Method != config.HTTPMethod {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if config.ContentType != "" {
		w.Header().Set("Content-Type", config.ContentType)
	}

	if config.CustomHeaders != "" {
		headers := strings.Split(config.CustomHeaders, "\n")
		for _, header := range headers {
			parts := strings.SplitN(header, ":", 2)
			if len(parts) == 2 {
				w.Header().Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
			}
		}
	}

	if config.FailureRate > 0 {
		if config.FailureRate > rand.Intn(100) {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(w, config.FailureResponseBody)
			return
		}
	}

	if config.HangForever {
		select {}
	}

	if config.DelaySeconds > 0 {
		select {
		case <-r.Context().Done():
			return
		case <-time.After(time.Duration(config.DelaySeconds) * time.Second):
		}
	}

	w.WriteHeader(config.StatusCode)
	_, _ = fmt.Fprint(w, config.ResponseBody)
}
