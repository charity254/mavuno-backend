package middleware

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the response writer to capture the status code
		wrapped := &responseWriter{w, http.StatusOK}

		next.ServeHTTP(wrapped, r)

		log.Printf(
			"%s %s %d %s",
			r.Method,
			r.RequestURI,
			wrapped.statusCode,
			time.Since(start),
		)
	})
}

// responseWriter wraps http.ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}