package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rhysmcneill/dividr/internal/config"
	mw "github.com/rhysmcneill/dividr/internal/middleware"
)

// TODO: Refactor handlers into their own packages/files, only have main.go handle server setup
func main() {
	// 1. Load Config
	cfg := config.Load()
	slog.Info("Configuration initialized", "config", cfg.LogValue())
	// 2. Init Logger
	slog.Info("Logger initialized", "level", cfg.LogLevel)

	// 3. Initialize Chi Router
	r := chi.NewRouter()

	// 4. Global Middleware (Applied to ALL requests)
	// Recoverer catches panics so your server doesn't crash
	r.Use(middleware.Recoverer)
	// Your Custom Logger (Injects Request ID + Slog)
	r.Use(mw.Logger)

	// 5. Define Routes (per your Roadmap Story 0.3.1)

	// Group: API (JSON only)
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			// This will now have a request_id logged automatically
			_, _ = w.Write([]byte(`{"status":"ok"}`))
		})

		r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
			// Placeholder for Epic 2.2
			w.WriteHeader(http.StatusUnauthorized)
		})
	})

	// Group: HTMX Fragments (HTML snippets)
	r.Route("/hx", func(r chi.Router) {
		r.Post("/validate-upload", func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("<div>Validating...</div>"))
		})
	})

	// Group: App Pages (The UI)
	r.Route("/app", func(r chi.Router) {
		// Middleware specific to the app (e.g., Auth Check)
		// r.Use(customMiddleware.RequireAuth)

		r.Get("/dashboard", func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("<h1>Dashboard Placeholder</h1>"))
		})
	})

	slog.Info("Starting Dividr HTTP server at", "port", cfg.Port)

	// 6. Start Server
	addr := fmt.Sprintf(":%s", cfg.Port)
	if err := http.ListenAndServe(addr, r); err != nil {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}
}
