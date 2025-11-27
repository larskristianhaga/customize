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
		DelaySeconds: 0,
		ResponseBody: `{"status": "success", "message": "This is a success example"}`,
		StatusCode:   200,
		HTTPMethod:   "GET",
	},
	"timeout": {
		DelaySeconds: 25,
		ResponseBody: `{"status": "success", "message": "This is a 25-second timeout example"}`,
		StatusCode:   200,
		HTTPMethod:   "GET",
	},
	"timeout-post": {
		DelaySeconds: 25,
		ResponseBody: `{"status": "success", "message": "This is a 25-second timeout example from POST"}`,
		StatusCode:   200,
		HTTPMethod:   "POST",
	},
	"created": {
		DelaySeconds:  0,
		ResponseBody:  `{"status": "created", "message": "Resource created successfully", "id": "42"}`,
		StatusCode:    201,
		HTTPMethod:    "POST",
		CustomHeaders: "Location: /api/v1/resources/42",
	},
	"badrequest": {
		DelaySeconds: 0,
		ResponseBody: `{"status": "error", "message": "Invalid request parameters"}`,
		StatusCode:   400,
		HTTPMethod:   "GET",
	},
	"unauthorized": {
		DelaySeconds: 0,
		ResponseBody: `{"status": "error", "message": "Unauthorized access"}`,
		StatusCode:   401,
		HTTPMethod:   "GET",
	},
	"forbidden": {
		DelaySeconds: 0,
		ResponseBody: `{"status": "error", "message": "Access forbidden"}`,
		StatusCode:   403,
		HTTPMethod:   "GET",
	},
	"notfound": {
		DelaySeconds: 0,
		ResponseBody: `{"status": "error", "message": "Resource not found"}`,
		StatusCode:   404,
		HTTPMethod:   "GET",
	},
	"ratelimit": {
		DelaySeconds: 0,
		ResponseBody: `{"status": "error", "message": "Too many requests"}`,
		StatusCode:   429,
		HTTPMethod:   "GET",
	},
	"teapot": {
		DelaySeconds: 0,
		ResponseBody: `{"status": "teapot", "message": "I'm a teapot"}`,
		StatusCode:   418,
		HTTPMethod:   "GET",
	},
	"error": {
		DelaySeconds: 0,
		ResponseBody: `{"status": "error", "message": "This is an error example"}`,
		StatusCode:   500,
		HTTPMethod:   "GET",
	},
	"serviceunavailable": {
		DelaySeconds: 0,
		ResponseBody: `{"status": "error", "message": "Service temporarily unavailable"}`,
		StatusCode:   503,
		HTTPMethod:   "GET",
	},
}
