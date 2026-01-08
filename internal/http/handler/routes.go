package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	mw "github.com/rhysmcneill/dividr/internal/http/middleware"
)

// RegisterRoutes sets up the router and all endpoints
func (h *Handler) RegisterRoutes() *chi.Mux {
	// Create a new router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Recoverer)
	r.Use(mw.Logger)

	// 1. Global Middleware
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	// r.Use(customMiddleware.Cors) - add later

	// 2. The Routes
	r.Get("/health", h.Health)

	// Future routes will go here:
	// r.Get("/login", h.LoginPage)
	// r.Post("/login", h.LoginSubmit)
	// r.Get("/dashboard", h.Dashboard)

	return r
}
