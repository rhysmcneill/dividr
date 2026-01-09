package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/rhysmcneill/dividr/internal/config"
	"github.com/rhysmcneill/dividr/internal/database"
	"github.com/rhysmcneill/dividr/internal/errs"
	"github.com/rhysmcneill/dividr/internal/http/json"
)

type Handler struct {
	DB     *database.Service
	Config *config.Config
}

func New(db *database.Service, cfg *config.Config) *Handler {
	return &Handler{
		DB:     db,
		Config: cfg,
	}
}

// respondWithError intelligently handles AppErrors vs System Errors
func (h *Handler) respondWithError(w http.ResponseWriter, r *http.Request, err error) {
	var appErr *errs.AppError

	// 1. Check if it is our custom "AppError"
	if errors.As(err, &appErr) {
		// Log based on the level defined in the error
		msg := "application error"
		attrs := []any{
			slog.String("code", appErr.Code),
			slog.String("path", r.URL.Path),
			slog.String("error", appErr.Err.Error()),
		}

		switch appErr.LogLevel {
		case slog.LevelError:
			slog.Error(msg, attrs...)
		case slog.LevelWarn:
			slog.Warn(msg, attrs...)
		default:
			slog.Info(msg, attrs...)
		}

		// Send JSON response and check for write errors
		if encErr := json.EncodeError(w, appErr.Status, appErr.Msg, appErr.Field); encErr != nil {
			slog.Error("failed to encode app error response", "error", encErr)
		}
		return
	}

	// 2. Unknown / System Error
	slog.Error("server error",
		"path", r.URL.Path,
		"error", err,
		"trace", string(debug.Stack()),
	)

	// Send 500 JSON and check for write errors
	if encErr := json.EncodeError(w, http.StatusInternalServerError, "The server encountered a problem", nil); encErr != nil {
		slog.Error("failed to encode 500 error response", "error", encErr)
	}
}
