package handlers

import (
	"net/http"
	"strings"
	"text/template"
)

func BasicApiExamplesHandler(w http.ResponseWriter, _ *http.Request) {
	tmpl, _ := template.ParseFiles("templates/basic-api-examples-documentation.html")
	_ = tmpl.Execute(w, nil)
}

func OpenAPIHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/yml")
	tmpl, _ := template.ParseFiles("files/openapi.yml")
	_ = tmpl.Execute(w, nil)
}

func ExamplesApiHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/")
	parts := strings.Split(path, "/")

	if len(parts) < 2 {
		http.Error(w, "Invalid API path", http.StatusBadRequest)
		return
	}

	exampleName := parts[1]
	if cfg, ok := exampleEndpoints[exampleName]; ok {
		HandleRequest(w, r, cfg)
		return
	}
	http.Error(w, "Example not found", http.StatusNotFound)
}

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
