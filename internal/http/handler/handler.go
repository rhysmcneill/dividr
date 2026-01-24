package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/rhysmcneill/dividr/internal/config"
	"github.com/rhysmcneill/dividr/internal/database"
	"github.com/rhysmcneill/dividr/internal/errs"
	"github.com/rhysmcneill/dividr/internal/http/handler/render"
	"github.com/rhysmcneill/dividr/internal/http/json"
	"github.com/rhysmcneill/dividr/web"
	"github.com/rhysmcneill/dividr/web/templates/dashboard"
	errpages "github.com/rhysmcneill/dividr/web/templates/errors"
	"github.com/rhysmcneill/dividr/web/templates/landing"
)

type Handler struct {
	DB             *database.Service
	Config         *config.Config
	SessionManager *scs.SessionManager
	RobotsData     []byte
	SitemapData    []byte
}

func New(db *database.Service, cfg *config.Config, sessionManager *scs.SessionManager) *Handler {
	// 1. Load SEO files
	robots, err := web.Files.ReadFile("seo/robots.txt")
	if err != nil {
		slog.Error("failed to load robots.txt", "error", err)
		panic("Fatal Error: robots.txt missing from embed")
	}

	sitemap, err := web.Files.ReadFile("seo/sitemap.xml")
	if err != nil {
		slog.Error("failed to load sitemap.xml", "error", err)
		panic("Fatal Error: sitemap.xml missing from embed")
	}

	return &Handler{
		DB:             db,
		Config:         cfg,
		SessionManager: sessionManager,
		RobotsData:     robots,
		SitemapData:    sitemap,
	}
}

// SEO HANDLERS
func (h *Handler) RobotsTXT(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write(h.RobotsData); err != nil {
		slog.Error("failed to write robots.txt", "error", err)
	}
}

func (h *Handler) SitemapXML(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write(h.SitemapData); err != nil {
		slog.Error("failed to write sitemap.xml", "error", err)
	}
}

// handleLanding renders the public home page
func (h *Handler) handleLanding(w http.ResponseWriter, r *http.Request) {
	// 200 OK
	if err := render.Component(w, r, http.StatusOK, landing.Page()); err != nil {
		h.respondWithError(w, r, err)
	}
}

// handleDashboard renders the main app view
func (h *Handler) handleDashboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. Get User ID from session
	// (The RequireAuth middleware guarantees this exists, so we can be confident)
	userIDStr := h.SessionManager.GetString(ctx, "userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		// Should rarely happen if middleware is working, but safety first
		slog.Error("invalid user_id in session", "error", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// 2. Check Database for HMRC Connection
	isConnected := false

	// Convert UUID to pgtype (assuming you use pgx/v5 and sqlc)
	// If you use a different driver, this might just be 'userID'
	conn, err := h.DB.GetHMRCConnectionByUserID(ctx, database.UUIDToPgtype(userID))

	// If err is nil, we found a record!
	if err == nil && conn.MtdID != "" {
		isConnected = true
	}
	// We ignore 'sql.ErrNoRows' because that just means isConnected = false (default)

	// 3. Render Dashboard with State
	// We pass the 'isConnected' flag to the template
	if err := render.Component(w, r, http.StatusOK, dashboard.Page(isConnected)); err != nil {
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
				if err := render.Component(w, r, appErr.Status, errpages.NotFound()); err != nil {
					slog.Error("failed to render 404 page", "error", err)
				}
			} else {
				// Otherwise show the generic 500 page with the specific error message
				if err := render.Component(w, r, appErr.Status, errpages.ServerError(appErr.Msg)); err != nil {
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
		if err := render.Component(w, r, http.StatusInternalServerError, errpages.ServerError("")); err != nil {
			slog.Error("failed to render 500 page", "error", err)
		}
	}
}
