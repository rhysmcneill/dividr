package response

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/rhysmcneill/dividr/internal/errs"
	mw "github.com/rhysmcneill/dividr/internal/middleware"
	"github.com/rhysmcneill/dividr/internal/ui/components"
)

// JSON Helper (Standard API)
func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode json", "error", err)
	}
}

// HTML Helper (For Templ Components)
// Usage: response.HTML(w, r, components.Dashboard(data))
func HTML(w http.ResponseWriter, r *http.Request, component templ.Component) {
	w.Header().Set("Content-Type", "text/html")
	// If it's a full page load, you might want to wrap it in a layout here,
	// but usually, that's handled at the handler level.
	_ = component.Render(r.Context(), w)
}

// HTMX Partial Helper
// Sometimes you want to set specific HTMX headers (like triggering a client-side event)
func HTMX(w http.ResponseWriter, r *http.Request, component templ.Component, headers map[string]string) {
	for k, v := range headers {
		w.Header().Set(k, v)
	}
	HTML(w, r, component)
}

// logError is a private helper that focuses ONLY on logging
func logError(r *http.Request, e *errs.AppError) {
	reqID, _ := r.Context().Value(mw.RequestIDKey).(string)

	// Create the logger with context
	log := slog.With(
		"request_id", reqID,
		"path", r.URL.Path,
		"status", e.Status,
		"code", e.Code,
	)

	// NEW: No switch statement. Just log at the level the error requested.
	switch e.LogLevel {
	case slog.LevelError:
		log.Error("request_failed", "err", e.Err)
	case slog.LevelWarn:
		log.Warn("request_warning")
	case slog.LevelInfo:
		log.Info("request_completed")
	case slog.LevelDebug:
		log.Debug("request_debug")
	default:
		// Fallback for safety
		log.Error("unknown_error_level", "err", e.Err)
	}
}

// 4. The Smart Error Handler
// This decides whether to send JSON or HTML based on the request headers.
func Error(w http.ResponseWriter, r *http.Request, err error) {
	// Convert generic error to AppError
	e, ok := err.(*errs.AppError)
	if !ok {
		e = errs.NewInternal(err)
	}

	// Log it (with Request ID for traceability)
	logError(r, e)

	// Content Negotiation: Is this an HTMX request?
	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(e.Status)
		// Render the type-safe Templ component
		_ = components.ErrorToast(e.Msg).Render(r.Context(), w)
		return
	}
	// Default to JSON
	JSON(w, e.Status, e)
}
