package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/rhysmcneill/dividr/internal/config"
)

// Init sets up the global logger.
// Call this once in main.go.
func Init(cfg *config.Config) {
	opts := &slog.HandlerOptions{
		// Add source file info only in dev (it's expensive in prod)
		//AddSource: cfg.AppEnv == "dev",
		// Set log level based on env
		Level: logLevel(cfg),
		// ReplaceAttr is used to scrub or rename fields
		ReplaceAttr: replaceAttr,
	}

	// Use JSON handler for production (structured)
	// You could use TextHandler for local dev if you prefer readability
	handler := slog.NewJSONHandler(os.Stdout, opts)

	// Create the logger
	logger := slog.New(handler)

	// Set as global default
	slog.SetDefault(logger)
}

func logLevel(cfg *config.Config) slog.Level {
	switch cfg.LogLevel {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN", "WARNING":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}

// replaceAttr is your "Sanitization" layer.
// This is where you ensure no Social Security numbers or sensitive tax data
// accidentally leaks into logs.
func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	// Standardize time format (currently using default format)
	// Uncomment and modify if custom time formatting is needed:
	// if a.Key == slog.TimeKey {
	// 	a.Value = slog.StringValue(a.Value.Time().Format(time.RFC3339))
	// }

	// Example: Mask generic keys if they ever slip in
	if a.Key == "password" || a.Key == "token" {
		return slog.String(a.Key, "[REDACTED]")
	}

	return a
}

// WithRequest returns a logger with the request_id attached.
// You use this in your handlers.
func WithRequest(ctx context.Context) *slog.Logger {
	// You will define your context keys in a separate 'keys' package to avoid cycles
	// For now, assuming you store request_id as a string in context
	reqID, _ := ctx.Value("request_id").(string)

	if reqID == "" {
		return slog.Default()
	}

	return slog.Default().With("request_id", reqID)
}
