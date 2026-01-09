package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	mw "github.com/rhysmcneill/dividr/internal/http/middleware"
)

// RegisterRoutes sets up the router and all endpoints
func (h *Handler) RegisterRoutes() *chi.Mux {
	r := chi.NewRouter()

	// Global Middleware (applied in order)
	r.Use(mw.RecoverPanic)      // Catch panics and log stack traces
	r.Use(middleware.RequestID) // Add chi's request ID to context
	r.Use(mw.Logger)            // Log requests with request ID correlation
	r.Use(mw.SecureHeaders)     // Add security headers

	// --- Public Routes ---
	r.HandleFunc("GET /health", h.handleDBHealth) // Basic connectivity check

	// --- API V1 Routes ---
	r.Route("/api/v1", func(r chi.Router) {

		// Public API Routes
		r.Get("/health", h.handleHealth)
		// r.Get("/login", h.LoginPage)
		// r.Post("/login", h.LoginSubmit)

		// // Private API routes
		r.Group(func(r chi.Router) {
			// r.Use(h.RequireToken) // <--- THIS is what makes it private later

			// r.HandleFunc("GET /", h.handleHome)
			// r.Get("/me", h.handleGetProfile)
			// r.Get("/dashboard", h.Dashboard)

		})
	})

	// The App (Strictly Private routes - Requires Session)
	r.Route("/app", func(r chi.Router) {
		// FUTURE: r.Use(AuthMiddleware)
		// FUTURE: r.Use(h.RequireSession)
		// TODO: r.Get("/dashboard",  h.handleDashboard)
	})

	// C. HX Group (Fragments Only - Requires Session Later)
	r.Route("/hx", func(r chi.Router) {
		// FUTURE: r.Use(AuthMiddleware)
		// FUTURE: r.Use(h.RequireSession)
		// TODO: r.Get("/summary", h.handleHXSummary)
	})

	return r
}
