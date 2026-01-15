package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rhysmcneill/dividr/internal/http/handler/render"
	mw "github.com/rhysmcneill/dividr/internal/http/middleware"
	"github.com/rhysmcneill/dividr/web"
	errpages "github.com/rhysmcneill/dividr/web/templates/errors"
)

// RegisterRoutes sets up the router and all endpoints
func (h *Handler) RegisterRoutes() *chi.Mux {
	r := chi.NewRouter()

	// Global Middleware (applied in order)
	r.Use(mw.RecoverPanic)      // Catch panics and log stack traces
	r.Use(middleware.RequestID) // Add chi's request ID to context
	r.Use(mw.Logger)            // Log requests with request ID correlation
	r.Use(mw.SecureHeaders)     // Add security headers
	r.Use(mw.CSPMiddleware)     // Content Security Policy headers with nonce

	// --- Static Files ---
	fs := http.FileServer(http.FS(web.Files))
	r.With(mw.CacheControl).Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	})

	// --- Error Pages ---
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		_ = render.Component(w, r, http.StatusNotFound, errpages.NotFound())
	})

	// --- Public Routes ---
	r.HandleFunc("GET /health", h.handleDBHealth) // Basic connectivity check
	r.Get("/", h.handleLanding)
	r.Post("/waitlist", WaitlistHandler(h))

	// Legal Pages
	r.Get("/privacy", PrivacyHandler)
	r.Get("/terms", TermsHandler)
	r.Get("/security", SecurityHandler)

	// Auth Routes
	r.Get("/login", h.handleLoginPage)
	r.Post("/login", h.handleLoginSubmit)

	r.Get("/signup", h.handleSignupPage)
	r.Post("/signup", h.handleSignupSubmit)

	r.Post("/logout", h.handleLogout)

	// --- API V1 Routes ---
	r.Route("/api/v1", func(r chi.Router) {

		// Public API Routes
		r.Get("/health", h.handleHealth)
		// r.Get("/login", h.LoginPage)
		// r.Post("/login", h.LoginSubmit)
		// Private API routes
		r.Group(func(r chi.Router) {
			// r.Use(h.RequireToken) // <--- THIS is what makes it private later

			// r.Get("/me", h.handleGetProfile)
			// r.Get("/dashboard", h.Dashboard)

		})
	})

	// The App (Strictly Private routes - Requires Session)
	r.Route("/app", func(r chi.Router) {
		r.Use(mw.RequireAuth(h.sessionManager))
		r.Get("/dashboard", h.handleDashboard)
	})

	// C. HX Group (Fragments Only - Requires Session Later)
	r.Route("/hx", func(r chi.Router) {
		r.Use(mw.RequireAuth(h.sessionManager))
		// TODO: r.Get("/summary", h.handleHXSummary)
	})

	return r
}
