package middleware

import (
	"net/http"
	"time"

	"calespero/pkg/logger"
)

// LoggingMiddleware wraps an http.Handler and logs request information
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a custom response writer to capture the status code
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Get user ID from context if available
		userID := "anonymous"
		if id, ok := r.Context().Value("user_id").(string); ok {
			userID = id
		}

		// Call the next handler
		next.ServeHTTP(rw, r)

		// Log the request
		logger.LogRequest(logger.RequestLog{
			Method:     r.Method,
			Path:       r.URL.Path,
			Duration:   time.Since(start),
			StatusCode: rw.statusCode,
			UserID:     userID,
		})
	}
}

// responseWriter is a custom response writer that captures the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
