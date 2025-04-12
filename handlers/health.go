package handlers

import (
	"net/http"
)

func HealthHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("I'm healthy"))
}
