package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"log"
	"time"	
	"strings"
	"sync"
	"math/rand"
	"github.com/google/uuid"
)

type UserConfig struct {
	DelaySeconds      int      `json:"delay_seconds"`
	ResponseBody      string   `json:"response_body"`
	StatusCode        int      `json:"status_code"`
	HangForever       bool     `json:"hang_forever"`
	HTTPMethod        string   `json:"http_method"`
	RandomDelay       bool     `json:"random_delay"`
	CustomHeaders     string   `json:"custom_headers"`
	FailureRate       int      `json:"failure_rate"`
	ResponseVariability string `json:"response_variability"`
}

var (
	configs = map[string]UserConfig{}
	mu      sync.RWMutex
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", landingHandler)
	http.HandleFunc("/dashboard", dashboardHandler)
	http.HandleFunc("/save", saveHandler)
	http.HandleFunc("/api/endpoint/", timeoutHandler)
	http.HandleFunc("/health", HealthHandler)

	log.Printf("🚀 Server started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func landingHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

func timeoutHandler(w http.ResponseWriter, r *http.Request) {
	userID := strings.TrimPrefix(r.URL.Path, "/api/endpoint/")
	
	mu.RLock()
	cfg, ok := configs[userID]
	mu.RUnlock()

	if !ok {
		http.Error(w, "User config not found", http.StatusNotFound)
		return
	}

	// Check if request method matches configured method
	if cfg.HTTPMethod != "" && r.Method != cfg.HTTPMethod {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
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
		status := http.StatusInternalServerError
		w.WriteHeader(status)
		fmt.Fprint(w, "Simulated failure")
		return
	}

	if cfg.HangForever {
		select {} // hang forever
	}

	// Calculate delay with variability if configured
	delay := time.Duration(cfg.DelaySeconds) * time.Second
	if cfg.RandomDelay {
		// Add random variation between -50% and +50%
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

func HealthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("I'm healthy"))
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Try to load template from standard location first
	tmpl, err := template.ParseFiles("/usr/local/share/customizeapi/templates/dashboard.html")
	if err != nil {
		// Fallback to local development path
		tmpl, err = template.ParseFiles("templates/dashboard.html")
		if err != nil {
			http.Error(w, "Error loading template: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	tmpl.Execute(w, nil)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userID := uuid.New().String()
	cfg := UserConfig{
		DelaySeconds:      atoi(r.FormValue("delay_seconds")),
		ResponseBody:      r.FormValue("response_body"),
		StatusCode:        atoi(r.FormValue("status_code")),
		HangForever:       r.FormValue("hang_forever") == "on",
		HTTPMethod:        r.FormValue("http_method"),
		RandomDelay:       r.FormValue("random_delay") == "on",
		CustomHeaders:     r.FormValue("custom_headers"),
		FailureRate:       atoi(r.FormValue("failure_rate")),
		ResponseVariability: r.FormValue("response_variability"),
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
	endpointURL := fmt.Sprintf("http://%s/api/endpoint/%s", host, userID)

	tmpl := template.Must(template.New("saved").Parse(`
		<div class="bg-green-50 text-green-800 rounded-lg p-4 mt-4">
			<div class="flex">
				<div class="flex-shrink-0">
					<svg class="h-5 w-5 text-green-400" viewBox="0 0 20 20" fill="currentColor">
						<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/>
					</svg>
				</div>
				<div class="ml-3">
					<p class="text-sm font-medium">
						Configuration saved successfully!
					</p>
					<div class="mt-4">
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
				</div>
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
	tmpl.Execute(w, data)
}

func atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}


