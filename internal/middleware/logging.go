package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// We use a custom type for context keys to prevent collisions
type contextKey string

const RequestIDKey contextKey = "request_id"

// Logger middleware
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 1. Check for existing correlation ID from headers (HMRC often sends one)
		// otherwise generate a new one.
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}

		// 2. Add request_id to headers so the client sees it (for support)
		w.Header().Set("X-Request-ID", reqID)

		// 3. Create a context with the request_id
		ctx := context.WithValue(r.Context(), RequestIDKey, reqID)

		// 4. Create a logger instance specifically for this request
		// This attaches the ID to every log line generated inside the handler
		log := slog.With(
			slog.String("request_id", reqID),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("ip", r.RemoteAddr), // Note: You'll need a real IP extractor later
		)

		// 5. Wrap the writer to capture the status code
		ww := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}

		// 6. Serve the request
		next.ServeHTTP(ww, r.WithContext(ctx))

		// 7. Log the outcome
		// Using generic "msg" makes it easy to search/filter later
		log.Info("http_request_completed",
			slog.Int("status", ww.statusCode),
			slog.Duration("duration", time.Since(start)),
			slog.String("user_agent", r.UserAgent()),
		)
	})
}

// Helper to capture status code
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriterWrapper) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}
