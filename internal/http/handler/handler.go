package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/rhysmcneill/dividr/internal/config"
	"github.com/rhysmcneill/dividr/internal/database"
	"github.com/rhysmcneill/dividr/internal/errs"
	"github.com/rhysmcneill/dividr/internal/http/handler/render" // <--- Need this
	"github.com/rhysmcneill/dividr/internal/http/json"
	"github.com/rhysmcneill/dividr/web/templates/pages" // <--- Need this
)

type Handler struct {
	DB     *database.Service
	Config *config.Config
}

func New(db *database.Service, cfg *config.Config) *Handler {
	return &Handler{
		DB:     db,
		Config: cfg,
	}
}

// handleLanding renders the public home page
func (h *Handler) handleLanding(w http.ResponseWriter, r *http.Request) {
	// 200 OK
	if err := render.Component(w, r, http.StatusOK, pages.Landing()); err != nil {
		h.respondWithError(w, r, err)
	}
}

// handleDashboard renders the main app view
func (h *Handler) handleDashboard(w http.ResponseWriter, r *http.Request) {
	// 200 OK
	if err := render.Component(w, r, http.StatusOK, pages.Dashboard()); err != nil {
		h.respondWithError(w, r, err)
	}
}

// respondWithError intelligently handles AppErrors vs System Errors
func (h *Handler) respondWithError(w http.ResponseWriter, r *http.Request, err error) {
	// 0. Detect "Zone" based on URL
	// If it starts with /api or /hx, we generally want JSON (or fragments).
	// For now, let's say /api gets JSON, everything else gets HTML pages.
	isAPI := strings.HasPrefix(r.URL.Path, "/api")

	var appErr *errs.AppError

	// --- CASE 1: Known Application Error ---
	if errors.As(err, &appErr) {
		// 1a. Log it (Reuse existing logging logic)
		msg := "application error"
		attrs := []any{
			slog.String("code", appErr.Code),
			slog.String("path", r.URL.Path),
			slog.String("error", appErr.Err.Error()),
		}

		switch appErr.LogLevel {
		case slog.LevelError:
			slog.Error(msg, attrs...)
		case slog.LevelWarn:
			slog.Warn(msg, attrs...)
		default:
			slog.Info(msg, attrs...)
		}

		// 1b. Send Response (Branch logic)
		if isAPI {
			if encErr := json.EncodeError(w, appErr.Status, appErr.Msg, appErr.Field); encErr != nil {
				slog.Error("failed to encode app error response", "error", encErr)
			}
		} else {
			// If it's a 404, show the specific Not Found page
			if appErr.Status == http.StatusNotFound {
				if err := render.Component(w, r, appErr.Status, pages.NotFound()); err != nil {
					slog.Error("failed to render 404 page", "error", err)
				}
			} else {
				// Otherwise show the generic 500 page with the specific error message
				if err := render.Component(w, r, appErr.Status, pages.ServerError(appErr.Msg)); err != nil {
					slog.Error("failed to render error page", "error", err)
				}
			}
		}
		return
	}

	// --- CASE 2: Unknown / System Error ---
	slog.Error("server error",
		"path", r.URL.Path,
		"error", err,
		"trace", string(debug.Stack()),
	)

	if isAPI {
		if encErr := json.EncodeError(w, http.StatusInternalServerError, "The server encountered a problem", nil); encErr != nil {
			slog.Error("failed to encode 500 error response", "error", encErr)
		}
	} else {
		// For browsers, show the generic "Something went wrong" page.
		// We pass an empty string to show the default friendly message defined in the template.
		if err := render.Component(w, r, http.StatusInternalServerError, pages.ServerError("")); err != nil {
			slog.Error("failed to render 500 page", "error", err)
		}
	}
}
