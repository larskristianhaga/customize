package main

import (
	"fmt"
	"github.com/google/uuid"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
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

var (
	configs = map[string]UserConfig{}
	mu      sync.RWMutex
	domain  = "https://customize.fly.dev"
)

var exampleEndpoints = map[string]UserConfig{
	"success": {
		DelaySeconds: 1,
		ResponseBody: `{"status": "success", "message": "This is a success example"}`,
		StatusCode:   200,
		HTTPMethod:   "GET",
	},
	"timeout": {
		DelaySeconds: 10,
		ResponseBody: `{"status": "success", "message": "This is a 10-second timeout example"}`,
		StatusCode:   200,
		HTTPMethod:   "GET",
	},
	"created": {
		DelaySeconds:  1,
		ResponseBody:  `{"status": "created", "message": "Resource created successfully", "id": "42"}`,
		StatusCode:    201,
		HTTPMethod:    "POST",
		CustomHeaders: "Location: /api/v1/resources/42",
	},
	"badrequest": {
		DelaySeconds: 1,
		ResponseBody: `{"status": "error", "message": "Invalid request parameters"}`,
		StatusCode:   400,
		HTTPMethod:   "GET",
	},
	"unauthorized": {
		DelaySeconds: 1,
		ResponseBody: `{"status": "error", "message": "Unauthorized access"}`,
		StatusCode:   401,
		HTTPMethod:   "GET",
	},
	"forbidden": {
		DelaySeconds: 1,
		ResponseBody: `{"status": "error", "message": "Access forbidden"}`,
		StatusCode:   403,
		HTTPMethod:   "GET",
	},
	"notfound": {
		DelaySeconds: 1,
		ResponseBody: `{"status": "error", "message": "Resource not found"}`,
		StatusCode:   404,
		HTTPMethod:   "GET",
	},
	"ratelimit": {
		DelaySeconds: 1,
		ResponseBody: `{"status": "error", "message": "Too many requests"}`,
		StatusCode:   429,
		HTTPMethod:   "GET",
	},
	"teapot": {
		DelaySeconds: 1,
		ResponseBody: `{"status": "teapot", "message": "I'm a teapot"}`,
		StatusCode:   418,
		HTTPMethod:   "GET",
	},
	"error": {
		DelaySeconds: 1,
		ResponseBody: `{"status": "error", "message": "This is an error example"}`,
		StatusCode:   500,
		HTTPMethod:   "GET",
	},
	"serviceunavailable": {
		DelaySeconds: 1,
		ResponseBody: `{"status": "error", "message": "Service temporarily unavailable"}`,
		StatusCode:   503,
		HTTPMethod:   "GET",
	},
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", LandingHandler)
	http.HandleFunc("/dashboard", DashboardHandler)
	http.HandleFunc("/save", SaveHandler)
	http.HandleFunc("/api/v1/", ApiHandler)
	http.HandleFunc("/examples", ExamplesHandler)
	http.HandleFunc("/api/v1/spec/openapi.yml", OpenAPIHandler)
	http.HandleFunc("/health", HealthHandler)
	http.HandleFunc("/robots.txt", RobotsHandler)
	http.HandleFunc("/sitemap.xml", SitemapHandler)

	log.Println("App live and listening on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func LandingHandler(w http.ResponseWriter, _ *http.Request) {
	tmpl, _ := template.ParseFiles("templates/landing.html")
	_ = tmpl.Execute(w, nil)
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Check for existing user ID cookie
	cookie, err := r.Cookie("user_id")
	var userID string
	var cfg UserConfig

	if err != nil || cookie.Value == "" {
		// Generate new user ID if none exists
		userID = uuid.New().String()
		// Set cookie with 1 year expiration
		http.SetCookie(w, &http.Cookie{
			Name:     "user_id",
			Value:    userID,
			Path:     "/",
			MaxAge:   365 * 24 * 60 * 60, // 1 year
			HttpOnly: true,
			Secure:   true,
		})

		// Initialize default config for new user
		cfg = UserConfig{
			DelaySeconds: 1,
			ResponseBody: "pong",
			StatusCode:   200,
			HTTPMethod:   "GET",
		}
		mu.Lock()
		configs[userID] = cfg
		mu.Unlock()
	} else {
		userID = cookie.Value
		mu.RLock()
		cfg = configs[userID]
		mu.RUnlock()
	}

	// Get the current host
	host := r.Host
	if host == "" {
		host = "localhost:8080"
	}

	// Create the full endpoint URL
	endpointURL := fmt.Sprintf("http://%s/api/v1/custom/%s", host, userID)

	tmpl, _ := template.ParseFiles("templates/dashboard.html")

	data := struct {
		EndpointURL      string
		Config          UserConfig
	}{
		EndpointURL: endpointURL,
		Config:      cfg,
	}
	tmpl.Execute(w, data)
}

func OpenAPIHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/yml")
	tmpl, _ := template.ParseFiles("files/openapi.yml")
	_ = tmpl.Execute(w, nil)
}

func ApiHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/")
	parts := strings.Split(path, "/")

	if len(parts) < 2 {
		http.Error(w, "Invalid API path", http.StatusBadRequest)
		return
	}

	// Handle example endpoints
	if parts[0] == "examples" && len(parts) == 2 {
		exampleName := parts[1]
		if cfg, ok := exampleEndpoints[exampleName]; ok {
			HandleRequest(w, r, cfg)
			return
		}
		http.Error(w, "Example not found", http.StatusNotFound)
		return
	}

	// Handle custom endpoints
	if parts[0] == "custom" && len(parts) == 2 {
		userID := parts[1]
		mu.RLock()
		cfg, ok := configs[userID]
		mu.RUnlock()

		if !ok {
			http.Error(w, "Endpoint not found", http.StatusNotFound)
			return
		}

		HandleRequest(w, r, cfg)
		return
	}

	http.Error(w, "Invalid API path", http.StatusBadRequest)
}

func HandleRequest(w http.ResponseWriter, r *http.Request, cfg UserConfig) {
	// Check if request method matches configured method
	if cfg.HTTPMethod != "" && r.Method != cfg.HTTPMethod {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set content type if configured
	if cfg.ContentType != "" {
		w.Header().Set("Content-Type", cfg.ContentType)
	}

	// Apply custom headers if configured
	if cfg.CustomHeaders != "" {
		headers := strings.Split(cfg.CustomHeaders, "\n")
		for _, header := range headers {
			parts := strings.SplitN(header, ":", 2)
			if len(parts) == 2 {
				w.Header().Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
			}
		}
	}

	// Check failure rate
	if cfg.FailureRate > 0 && rand.Intn(100) < cfg.FailureRate {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, cfg.FailureResponseBody)
		return
	}

	if cfg.HangForever {
		select {} // hang forever
	}

	// Calculate delay with variability if configured
	delay := time.Duration(cfg.DelaySeconds) * time.Second
	if cfg.RandomDelay {
		variation := rand.Float64() - 0.5
		delay = time.Duration(float64(delay) * (1 + variation))
	} else if cfg.ResponseVariability != "none" {
		var multiplier float64
		switch cfg.ResponseVariability {
		case "low":
			multiplier = 0.1
		case "medium":
			multiplier = 0.25
		case "high":
			multiplier = 0.5
		}
		variation := (rand.Float64()*2 - 1) * multiplier
		delay = time.Duration(float64(delay) * (1 + variation))
	}

	time.Sleep(delay)
	w.WriteHeader(cfg.StatusCode)
	fmt.Fprint(w, cfg.ResponseBody)
}

func SaveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	
	// Get user ID from cookie
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
	}
	
	mu.Lock()
	configs[userID] = cfg
	mu.Unlock()

	// Get the current host
	host := r.Host
	if host == "" {
		host = "localhost:8080"
	}

	// Create the full endpoint URL
	endpointURL := fmt.Sprintf("http://%s/api/v1/custom/%s", host, userID)

	tmpl := template.Must(template.New("saved").Parse(`
		<div class="bg-gray-50 p-4 rounded-lg">
			<p class="text-sm font-medium text-gray-700">Your API Endpoint:</p>
			<div class="mt-2 bg-gray-100 p-3 rounded-lg">
				<code class="text-blue-600 break-all">{{.EndpointURL}}</code>
			</div>
			<div class="mt-2 text-sm text-gray-600">
				<p>Example usage:</p>
				<pre class="mt-1 bg-gray-100 p-2 rounded text-sm">
curl -X {{.Method}} {{.EndpointURL}}
				</pre>
			</div>
		</div>
	`))

	data := struct {
		EndpointURL string
		Method      string
	}{
		EndpointURL: endpointURL,
		Method:      cfg.HTTPMethod,
	}

	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, data)
}

func atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func HealthHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("I'm healthy"))
}

func RobotsHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	_, _ = fmt.Fprint(w, `User-agent: *
Allow: /
Allow: /dashboard
Disallow: /api/v1/
Disallow: /api/v1/examples/
Disallow: /api/v1/custom/
Disallow: /save

Sitemap: `+domain+`/sitemap.xml`)
}

func SitemapHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/xml")
	_, _ = fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
    <url>
        <loc>`+domain+`</loc>
    </url>
    <url>
        <loc>`+domain+`/dashboard</loc>
    </url>
</urlset>`)
}

func ExamplesHandler(w http.ResponseWriter, _ *http.Request) {
	tmpl, _ := template.ParseFiles("templates/basic-api-examples-documentation.html")
	_ = tmpl.Execute(w, nil)
}
