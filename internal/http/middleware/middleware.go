package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/alexedwards/scs/v2"
)

// RecoverPanic catches any code that crashes the thread and returns a 500
func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the stack trace so we can debug it
				slog.Error("panic recovered",
					"error", err,
					"stack", string(debug.Stack()),
				)

				// Return a generic 500 error
				w.Header().Set("Connection", "close")
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// CacheControl sets the Cache-Control header.
// Used for static assets.
func CacheControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set aggressive caching (1 year)
		w.Header().Set("Cache-Control", "public, max-age=31536000")
		next.ServeHTTP(w, r)
	})
}

// SecureHeaders adds security headers to every response
func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

// RequireAuth redirects unauthenticated users to the login page
func RequireAuth(sessionManager *scs.SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// We check if "userID" exists in the session.
			// (We haven't built the Login handler yet, but this is the check we will use)
			if !sessionManager.Exists(r.Context(), "userID") {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
