package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/alexedwards/scs/v2"
)

const NonceKey contextKey = "cspNonce"

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

func SEOCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=86400")
		next.ServeHTTP(w, r)
	})
}

// ContentTypeText sets the Content-Type to text/plain
func ContentTypeText(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

// ContentTypeXML sets the Content-Type to application/xml
func ContentTypeXML(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

// CSPMiddleware generates a unique "Nonce" (token) for every request.
// It allows HTMX inline styles to work securely.
func CSPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Generate a random 16-byte hex string
		bytes := make([]byte, 16)
		if _, err := rand.Read(bytes); err != nil {
			slog.Error("failed to generate csp nonce", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		nonce := hex.EncodeToString(bytes)

		// 2. Build the dynamic CSP header
		// We tell the browser: "Allow scripts and styles ONLY if they match this nonce"
		csp := fmt.Sprintf(
			"default-src 'self'; "+
				"style-src 'self' 'nonce-%s' fonts.googleapis.com; "+
				"script-src 'self' 'nonce-%s'; "+
				"font-src 'self' fonts.gstatic.com",
			nonce, nonce,
		)

		w.Header().Set("Content-Security-Policy", csp)

		// 3. Store the nonce in the context so your Templ files can read it
		ctx := context.WithValue(r.Context(), NonceKey, nonce)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetNonce is a helper function your Templates will call to get the token.
func GetNonce(ctx context.Context) string {
	if nonce, ok := ctx.Value(NonceKey).(string); ok {
		return nonce
	}
	return ""
}

// GetHTMXConfig returns the HTMX config meta content with the CSP nonce
func GetHTMXConfig(ctx context.Context) string {
	return fmt.Sprintf(`{"inlineStyleNonce":"%s"}`, GetNonce(ctx))
}
