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
}

func HandleAPIRequest(w http.ResponseWriter, r *http.Request, cfg UserConfig) {
	if cfg.HTTPMethod != "" && r.Method != cfg.HTTPMethod {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
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
