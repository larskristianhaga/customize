package middleware

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
)

// LoggingMiddleware logs the IP address and request details
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the IP address
		ip := r.RemoteAddr
		if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
			ip = ip
		}

		// Format the log message
		logMsg := fmt.Sprintf("[%s] %s %s %s\n",
			time.Now().Format("2006-01-02 15:04:05"),
			ip,
			r.Method,
			r.URL.Path,
		)

		// Write to stdout
		fmt.Fprint(os.Stdout, logMsg)

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
