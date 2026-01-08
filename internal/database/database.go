package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rhysmcneill/dividr/internal/config"
)

// Service wraps the SQLC Queries and the Connection Pool
type Service struct {
	*Queries // Embeds the SQLC methods (GetUsers, etc.)
	Pool     *pgxpool.Pool
}

// Connect initializes the DB connection and returns the Service
// NOTE: We renamed this from 'New' to 'Connect' to avoid conflict with SQLC's New()
func Connect(cfg *config.Config) (*Service, error) {
	// 1. Create the Connection Pool (Thread-safe, high performance)
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		slog.Error("Failed to create connection pool", "error", err)
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	// 2. Fast Ping to verify
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		slog.Error("Database ping failed", "error", err)
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	// 3. Initialize SQLC with the pool
	// 'New' here comes from the generated code in db.go
	queries := New(pool)

	slog.Info("Database connected successfully")

	return &Service{
		Queries: queries,
		Pool:    pool,
	}, nil
}

// Close closes the connection pool
func (s *Service) Close() {
	s.Pool.Close()
}

// Health checks the DB status (Re-adding this so your /health endpoint works)
func (s *Service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	status := "up"
	if err := s.Pool.Ping(ctx); err != nil {
		status = "down"
	}

	return map[string]string{
		"status":  status,
		"db_type": "postgres",
	}
}
