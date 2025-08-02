package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"text/template"

	"github.com/google/uuid"
)

var (
	configs = map[string]UserConfig{}
	mu      sync.RWMutex
)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Check for existing user ID cookie
	cookie, err := r.Cookie("user_id")
	var userID string
	var cfg UserConfig

	if err != nil || cookie.Value == "" {
		userID = uuid.New().String()
		http.SetCookie(w, &http.Cookie{
			Name:     "user_id",
			Value:    userID,
			Path:     "/",
			MaxAge:   5 * 60,
			HttpOnly: true,
			Secure:   true,
		})

		// Initialize default config for new user
		cfg = UserConfig{
			DelaySeconds: 1,
			ContentType:  "text/plain",
			ResponseBody: "pong",
			StatusCode:   200,
			HTTPMethod:   "GET",
			// Default CORS configuration
			EnableCORS:       false,
			CORSAllowOrigin:  "*",
			CORSAllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
			CORSAllowHeaders: "Content-Type, Authorization",
			CORSMaxAge:       86400, // 24 hours
			CORSAllowCreds:   false,
		}
		mu.Lock()
		configs[userID] = cfg
		mu.Unlock()
	} else {
		userID = cookie.Value
		mu.RLock()
		var ok bool
		cfg, ok = configs[userID]
		mu.RUnlock()

		// If config doesn't exist, create default one
		if !ok {
			cfg = UserConfig{
				DelaySeconds: 1,
				ContentType:  "text/plain",
				ResponseBody: "pong",
				StatusCode:   200,
				HTTPMethod:   "GET",
				// Default CORS configuration
				EnableCORS:       false,
				CORSAllowOrigin:  "*",
				CORSAllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
				CORSAllowHeaders: "Content-Type, Authorization",
				CORSMaxAge:       86400, // 24 hours
				CORSAllowCreds:   false,
			}
			mu.Lock()
			configs[userID] = cfg
			mu.Unlock()
		}
	}

	endpointURL := fmt.Sprintf("http://%s/api/v1/custom/%s", r.Host, userID)
	tmpl, _ := template.ParseFiles("templates/dashboard.html")

	data := struct {
		EndpointURL string
		Config      UserConfig
	}{
		EndpointURL: endpointURL,
		Config:      cfg,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("Error: ", err)
	}
}

func SaveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()

	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Error(w, "User ID not found", http.StatusBadRequest)
		return
	}
	userID := cookie.Value

	cfg := UserConfig{
		DelaySeconds:        atoi(r.FormValue("delay_seconds")),
		ResponseBody:        r.FormValue("response_body"),
		StatusCode:          atoi(r.FormValue("status_code")),
		HangForever:         r.FormValue("hang_forever") == "on",
		HTTPMethod:          r.FormValue("http_method"),
		RandomDelay:         r.FormValue("random_delay") == "on",
		CustomHeaders:       r.FormValue("custom_headers"),
		FailureRate:         atoi(r.FormValue("failure_rate")),
		ResponseVariability: r.FormValue("response_variability"),
		ContentType:         r.FormValue("content_type"),
		FailureResponseBody: r.FormValue("failure_response_body"),
		// CORS configuration
		EnableCORS:       r.FormValue("enable_cors") == "on",
		CORSAllowOrigin:  r.FormValue("cors_allow_origin"),
		CORSAllowMethods: r.FormValue("cors_allow_methods"),
		CORSAllowHeaders: r.FormValue("cors_allow_headers"),
		CORSMaxAge:       atoi(r.FormValue("cors_max_age")),
		CORSAllowCreds:   r.FormValue("cors_allow_credentials") == "on",
	}

	mu.Lock()
	configs[userID] = cfg
	mu.Unlock()

	w.WriteHeader(http.StatusOK)
}

func atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func CustomApiHandler(w http.ResponseWriter, r *http.Request) {
	// Remove the prefix from the path
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/custom/")
	// Split the remaining path
	parts := strings.Split(path, "/")

	userID := parts[0]

	mu.RLock()
	cfg, ok := configs[userID]
	mu.RUnlock()

	if !ok {
		http.Error(w, "Endpoint not found", http.StatusNotFound)
		return
	}

	HandleAPIRequest(w, r, cfg)
}
