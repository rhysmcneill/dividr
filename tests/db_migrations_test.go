package tests

import (
	"database/sql"
	"testing"

	"github.com/rhysmcneill/dividr/internal/config"
	"github.com/rhysmcneill/dividr/internal/database"

	// We need the stdlib driver for the test connection
	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestDatabaseMigrations(t *testing.T) {
	// 1. Setup Connection
	cfg := config.Load()
	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		t.Fatalf("Failed to open db connection: %v", err)
	}
	defer db.Close()

	// --- THE FIX ---
	// Drop the hidden migration history table too!
	// This forces the library to re-run your .sql files.
	_, err = db.Exec(`
        DROP TABLE IF EXISTS sessions;
        DROP TABLE IF EXISTS schema_migrations;
    `)
	if err != nil {
		t.Fatalf("Failed to clean database: %v", err)
	}

	// 3. ACT: Run the migration logic
	if err := database.RunMigrations(db); err != nil {
		t.Fatalf("RunMigrations failed: %v", err)
	}

	// 4. ASSERT: Check if the table exists now
	_, err = db.Exec("SELECT token, data, expiry FROM sessions LIMIT 1")
	if err != nil {
		t.Fatalf("Migration failed to create 'sessions' table: %v", err)
	}

	t.Log("Success: Migration ran and 'sessions' table exists.")
}
