package handler

import (
	"net/http"

	// Make sure this path matches your module name exactly
	"github.com/rhysmcneill/dividr/web/templates/legal"
)

// PrivacyHandler serves the privacy policy
func PrivacyHandler(w http.ResponseWriter, r *http.Request) {
	// It's good practice to explicitly set the content type for HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Render the component
	// We use r.Context() to ensure if the request is cancelled, rendering stops
	err := legal.Privacy().Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}

// TermsHandler serves the terms of service
func TermsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := legal.Terms().Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}

// SecurityHandler serves the security policy
func SecurityHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := legal.Security().Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}
