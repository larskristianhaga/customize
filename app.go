package main

import (
	"log"
	"net/http"
	"os"

	"github.com/larskristianhaga/customize/handlers"
	"github.com/larskristianhaga/customize/middleware"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register all handlers
	mux.HandleFunc("/", handlers.LandingHandler)
	mux.HandleFunc("/dashboard", handlers.DashboardHandler)
	mux.HandleFunc("/save", handlers.SaveHandler)
	mux.HandleFunc("/api/v1/custom/", handlers.CustomApiHandler)
	mux.HandleFunc("/basic-api-examples", handlers.BasicApiExamplesHandler)
	mux.HandleFunc("/api/v1/spec/openapi.yml", handlers.OpenAPIHandler)
	mux.HandleFunc("/api/v1/examples/", handlers.ExamplesApiHandler)
	mux.HandleFunc("/health", handlers.HealthHandler)
	mux.HandleFunc("/robots.txt", handlers.RobotsHandler)
	mux.HandleFunc("/sitemap.xml", handlers.SitemapHandler)

	// Wrap the mux with middleware
	handler := middleware.LoggingMiddleware(mux)
	handler = middleware.CORSMiddleware(handler)

	log.Println("App live and listening on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
