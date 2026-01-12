package handler

import (
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/rhysmcneill/dividr/internal/database"
	"github.com/rhysmcneill/dividr/web/templates/components"
	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(rate.Every(3*time.Second), 3)

func WaitlistHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// 1. In-Memory Rate Limit
		if !limiter.Allow() {
			// Render error component instead of plain text
			renderError(w, r, "Too many requests. Try again later.")
			return
		}

		if err := r.ParseForm(); err != nil {
			renderError(w, r, "Unable to process form.")
			return
		}

		// 2. Honeypot Check
		if r.FormValue("website_hp") != "" {
			// Bot detected: Show success to fool them
			renderSuccess(w, r)
			return
		}

		email := strings.TrimSpace(r.FormValue("email"))

		// 3. Validate Email
		if _, err := mail.ParseAddress(email); err != nil {
			renderError(w, r, "Invalid email address.")
			return
		}

		// 4. Save to DB
		queries := database.New(h.DB.Pool)
		_, err := queries.CreateWaitlistEntry(r.Context(), email)

		if err != nil {
			// Duplicate key error? Show success anyway (Privacy)
			if strings.Contains(err.Error(), "duplicate key value") {
				renderSuccess(w, r)
				return
			}
			// Real DB error
			renderError(w, r, "Something went wrong. Please try again.")
			return
		}

		// 5. Success
		renderSuccess(w, r)
	}
}

// Helper to render success
func renderSuccess(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := components.WaitlistSuccess().Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Rendering error", http.StatusInternalServerError)
	}
}

// Helper to render error
func renderError(w http.ResponseWriter, r *http.Request, msg string) {
	w.Header().Set("Content-Type", "text/html")
	// Note: We return 200 OK here so HTMX swaps the HTML.
	// If we returned 400/500, HTMX would ignore the response content by default.
	err := components.WaitlistError(msg).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Rendering error", http.StatusInternalServerError)
	}
}
