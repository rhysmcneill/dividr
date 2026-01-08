package handler

import (
	"net/http"

	"github.com/rhysmcneill/dividr/internal/database"
	"github.com/rhysmcneill/dividr/internal/http/response" // Import the toolkit
)

// The Handler struct holds your dependencies
type Handler struct {
	DB *database.Service
}

func New(db *database.Service) *Handler {
	return &Handler{DB: db}
}

// The methods use the response toolkit
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	// 1. Do logic using h.DB
	data := h.DB.Health()

	// 2. Use the toolkit to reply
	response.JSON(w, http.StatusOK, data)
}
