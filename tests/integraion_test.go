package tests

import (
	"context"
	"testing"
	"time"

	"github.com/rhysmcneill/dividr/internal/config"
	"github.com/rhysmcneill/dividr/internal/database"
)

// Correct signature: Takes *testing.T, returns nothing
func TestDatabaseConnection(t *testing.T) {
	// 1. Load config
	cfg := config.Load()

	// 2. Attempt connection using your app's actual logic
	// We assume database.Connect returns the pool/db object
	db, err := database.Connect(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	// Always close the pool at the end of the test
	defer db.Close()

	// 3. Verify it's alive
	// Create a context with a timeout so the test doesn't hang forever if DB is down
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := db.Pool.Ping(ctx); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}
}
