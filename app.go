package main

import (
	"log"
	"net/http"
	"os"

	"github.com/larskristianhaga/customize/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", handlers.LandingHandler)

	http.HandleFunc("/dashboard", handlers.DashboardHandler)
	http.HandleFunc("/save", handlers.SaveHandler)
	http.HandleFunc("/api/v1/custom/", handlers.CustomApiHandler)

	http.HandleFunc("/basic-api-examples", handlers.BasicApiExamplesHandler)
	http.HandleFunc("/api/v1/spec/openapi.yml", handlers.OpenAPIHandler)
	http.HandleFunc("/api/v1/examples/", handlers.ExamplesApiHandler)

	http.HandleFunc("/health", handlers.HealthHandler)

	http.HandleFunc("/robots.txt", handlers.RobotsHandler)
	http.HandleFunc("/sitemap.xml", handlers.SitemapHandler)

	log.Println("App live and listening on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
