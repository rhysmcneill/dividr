package tests

import (
	"context"
	"os"
	"strings"
	"testing"

	// Added missing import
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/rhysmcneill/dividr/internal/config"
	"github.com/rhysmcneill/dividr/internal/database" // Using your existing package
	"github.com/rhysmcneill/dividr/internal/hmrc"
)

func TestSimulationMode(t *testing.T) {
	// 1. Load .env
	if err := godotenv.Load("../.env"); err != nil {
		_ = godotenv.Load("../../.env") // Fallback for nested tests
	}

	// 2. FIX: Manually read the secret key if env var is missing
	// Because config.go looks in "./docker/...", which fails from the "tests" folder.
	if os.Getenv("TOKEN_ENCRYPTION_KEY") == "" {
		keyPath := "../docker/secrets/.crypto_key"

		if _, err := os.Stat(keyPath); os.IsNotExist(err) {
			keyPath = "./docker/secrets/.crypto_key" // Fallback for nested tests
		}

		key, err := os.ReadFile(keyPath)
		if err != nil {
			t.Logf("⚠️ Could not find secret at %s, and .env failed.", keyPath)
			t.Skip("Skipping test: Missing Encryption Key")
		}
		// Set it manually for this test process
		t.Setenv("TOKEN_ENCRYPTION_KEY", strings.TrimSpace(string(key)))
	}

	// 2. Load Config & Connect to DB
	// We use your 'database.Connect' directly now
	cfg := config.Load()
	dbService, err := database.Connect(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	defer dbService.Close()

	// 3. Define the Test User
	targetUserID := "9b21da56-fba5-4118-a954-097b5f9db9d1"
	userID, err := uuid.Parse(targetUserID)
	if err != nil {
		t.Fatalf("Invalid User UUID in test setup: %v", err)
	}

	ctx := context.Background()

	// 4. Fetch Credentials
	// dbService embeds *Queries, so we call methods directly on it
	t.Logf("🔍 Fetching credentials for user %s...", targetUserID)
	conn, err := dbService.GetHMRCConnectionByUserID(ctx, database.UUIDToPgtype(userID))
	if err != nil {
		t.Fatalf("No HMRC connection found for test user. Log in via web first.")
	}

	// 5. Decrypt Token
	accessToken, err := hmrc.Decrypt(conn.AccessToken, cfg.TokenEncryptionKey)
	if err != nil {
		t.Fatalf("Token Decryption failed. Check TOKEN_ENCRYPTION_KEY. Error: %v", err)
	}
	t.Logf("✅ Token Decrypted Successfully")

	// 6. Build Client
	meta := hmrc.FraudMetadata{
		ClientIP:        "127.0.0.1",
		ClientPort:      "8080",
		ClientDeviceID:  "test-suite-sim-" + uuid.New().String(),
		ClientUserAgent: "Dividr-Test-Suite/1.0",
		ClientTimezone:  "UTC",
		ClientUserID:    targetUserID,
	}

	client := hmrc.NewClient(false, accessToken, meta, cfg)

	// 7. Execute Request
	t.Log("🚀 Sending Request to HMRC Sandbox...")
	_, err = client.GetObligations(conn.MtdID, "2025-04-06", "2026-04-05")
	if err != nil {
		t.Fatalf("❌ HMRC Call Failed: %v", err)
	}

	t.Log("✅ SUCCESS: HMRC accepted Auth, Headers, and Encryption.")
}
