package handler

import (
	"net/http"

	"github.com/rhysmcneill/dividr/internal/http/json"
	// Import the toolkit
)

// Version can be injected at build time
var Version = "X.Y.Z"

// The method is used for public health checks
func (h *Handler) handleDBHealth(w http.ResponseWriter, r *http.Request) {
	data := h.DB.Health()
	// Use the json helper we created
	if err := json.Encode(w, http.StatusOK, data); err != nil {
		h.respondWithError(w, r, err)
	}
}

// The method is used for authenticated health checks
func (h *Handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	// 1. Get DB Health
	dbHealth := h.DB.Health() // Returns map[string]string

	// 2. Determine overall system status
	systemStatus := "available"
	statusCode := http.StatusOK

	if dbHealth["status"] != "up" {
		systemStatus = "unavailable"
		statusCode = http.StatusServiceUnavailable
	}

	// 3. Construct Response
	response := map[string]any{
		"status":      systemStatus,
		"environment": h.Config.AppEnv,
		"version":     Version,
		"components": map[string]any{
			"database": dbHealth,
		},
	}

	// 4. Send
	if err := json.Encode(w, statusCode, response); err != nil {
		h.respondWithError(w, r, err)
	}
}
