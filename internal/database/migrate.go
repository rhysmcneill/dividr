package database

import (
	"database/sql"
	"embed"
	"errors"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// RunMigrations automatically applies all up migrations embedded in the binary.
func RunMigrations(db *sql.DB) error {
	slog.Info("Checking for migrations...")

	// 1. Create a driver for the embedded filesystem (our SQL files)
	driver, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		slog.Error("failed to create iofs driver", "error", err)
		return err
	}

	// 2. Create the postgres driver
	dbDriver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		slog.Error("failed to create postgres driver", "error", err)
		return err
	}

	// 3. Create the migrate instance
	m, err := migrate.NewWithInstance(
		"iofs", driver, // Source: Embedded files
		"postgres", dbDriver, // Target: Connected DB
	)
	if err != nil {
		slog.Error("failed to create migration instance", "error", err)
		return err
	}

	// 4. Run "Up" (apply all changes)
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			slog.Info("Database Schema is up to date.")
			return nil
		}
		slog.Error("failed to apply migrations", "error", err)
		return err
	}

	slog.Info("Migrations applied successfully.")
	return nil
}
