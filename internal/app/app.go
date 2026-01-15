package app

import (
	"context"
	"database/sql" // <--- Added
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexedwards/scs/postgresstore" // <--- Added
	"github.com/alexedwards/scs/v2"            // <--- Added
	_ "github.com/jackc/pgx/v5/stdlib"         // <--- Added for SCS

	"github.com/rhysmcneill/dividr/internal/config"
	"github.com/rhysmcneill/dividr/internal/database"
	"github.com/rhysmcneill/dividr/internal/http/handler"
	"github.com/rhysmcneill/dividr/internal/logger"
)

func Run() error {
	cfg := config.Load()
	logger.Init(cfg.LogLevel)
	slog.Info("Configuration initialized", "config", cfg.LogValue())

	// 1. Main App Database (PGX Pool)
	db, err := database.Connect(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	// 2. Session Database (Standard SQL)
	// SCS requires a *sql.DB. We use the same DSN from config.
	sessionDB, err := sql.Open("pgx", cfg.DatabaseURL) // Ensure cfg has DatabaseURL or construct it
	if err != nil {
		return fmt.Errorf("failed to open session db: %w", err)
	}
	defer func() {
		if err := sessionDB.Close(); err != nil {
			slog.Error("failed to close session db", "error", err)
		}
	}()

	if err := database.RunMigrations(sessionDB); err != nil {
		slog.Error("failed to run session db migrations", "error", err)
		return err
	}

	// 3. Initialize Session Manager
	sm := scs.New()
	sm.Lifetime = 24 * time.Hour
	sm.Store = postgresstore.New(sessionDB)
	sm.Cookie.Name = "dividr_session"
	sm.Cookie.HttpOnly = true
	sm.Cookie.Secure = cfg.AppEnv == "production" // Secure in prod
	sm.Cookie.SameSite = http.SameSiteLaxMode

	// 4. Inject Session Manager into Handler
	// You must update handler.New to accept this!
	h := handler.New(db, cfg, sm)

	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", cfg.Port),
		// 5. Wrap Routes with Session Middleware
		Handler:      sm.LoadAndSave(h.RegisterRoutes()),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		slog.Info("Starting server", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Fatal Server error", "error", err)
			quit <- syscall.SIGTERM
		}
	}()

	// Block until signal received
	<-quit
	slog.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	slog.Info("Server stopped gracefully")
	return nil
}
