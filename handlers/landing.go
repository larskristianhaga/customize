package handlers

import (
	"net/http"
	"text/template"
)

func LandingHandler(w http.ResponseWriter, _ *http.Request) {
	tmpl, _ := template.ParseFiles("templates/landing.html")
	_ = tmpl.Execute(w, nil)
}
