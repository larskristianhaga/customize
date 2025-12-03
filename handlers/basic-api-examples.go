package handlers

import (
	"net/http"
	"strings"
	"text/template"
)

func BasicApiExamplesHandler(w http.ResponseWriter, _ *http.Request) {
	tmpl, _ := template.ParseFiles("templates/basic-api-examples.html")
	_ = tmpl.Execute(w, nil)
}

func OpenAPIHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/yml")
	tmpl, _ := template.ParseFiles("files/openapi.yml")
	_ = tmpl.Execute(w, nil)
}

func ExamplesApiHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	exampleName := parts[len(parts)-1]

	if cfg, ok := exampleEndpoints[exampleName]; ok {
		HandleAPIRequest(w, r, cfg)
		return
	}

	http.Error(w, "Example API not found", http.StatusNotFound)
}

var exampleEndpoints = map[string]UserConfig{
	"success": {
		HTTPMethod:   "GET",
		DelaySeconds: 0,
		StatusCode:   200,
		ResponseBody: `{"status": "success", "message": "This is a success example"}`,
		ContentType:  "application/json",
	},
	"created": {
		HTTPMethod:   "POST",
		DelaySeconds: 0,
		StatusCode:   201,
		ResponseBody: `{"status": "created", "message": "Resource created successfully", "id": "42"}`,
		ContentType:  "application/json",
	},
	"bad-request": {
		HTTPMethod:   "GET",
		DelaySeconds: 0,
		StatusCode:   400,
		ResponseBody: `{"status": "error", "message": "Invalid request parameters"}`,
		ContentType:  "application/json",
	},
	"unauthorized": {
		HTTPMethod:   "GET",
		DelaySeconds: 0,
		StatusCode:   401,
		ResponseBody: `{"status": "error", "message": "Unauthorized access"}`,
		ContentType:  "application/json",
	},
	"forbidden": {
		HTTPMethod:   "GET",
		DelaySeconds: 0,
		StatusCode:   403,
		ResponseBody: `{"status": "error", "message": "Access forbidden"}`,
		ContentType:  "application/json",
	},
	"not-found": {
		HTTPMethod:   "GET",
		DelaySeconds: 0,
		StatusCode:   404,
		ResponseBody: `{"status": "error", "message": "Resource not found"}`,
		ContentType:  "application/json",
	},
	"rate-limit": {
		HTTPMethod:   "GET",
		DelaySeconds: 0,
		StatusCode:   429,
		ResponseBody: `{"status": "error", "message": "Too many requests"}`,
		ContentType:  "application/json",
	},
	"teapot": {
		HTTPMethod:   "GET",
		DelaySeconds: 0,
		StatusCode:   418,
		ResponseBody: `{"status": "teapot", "message": "I'm a teapot"}`,
		ContentType:  "application/json",
	},
	"error": {
		HTTPMethod:   "GET",
		DelaySeconds: 0,
		StatusCode:   500,
		ResponseBody: `{"status": "error", "message": "This is an error example"}`,
		ContentType:  "application/json",
	},
	"error-half": {
		HTTPMethod:   "GET",
		DelaySeconds: 0,
		FailureRate:  50,
		StatusCode:   200,
		ResponseBody: `{"status": "error", "message": "This is an error example"}`,
		ContentType:  "application/json",
	},
	"service-unavailable": {
		HTTPMethod:   "GET",
		DelaySeconds: 0,
		StatusCode:   503,
		ResponseBody: `{"status": "error", "message": "Service temporarily unavailable"}`,
		ContentType:  "application/json",
	},
	"timeout": {
		HTTPMethod:   "GET",
		DelaySeconds: 30,
		StatusCode:   200,
		ResponseBody: `{"status": "success", "message": "This is a 30-second timeout example"}`,
		ContentType:  "application/json",
	},
}
