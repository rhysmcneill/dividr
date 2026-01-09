package json

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Response is the standard envelope for all API responses
type Response struct {
	Data  any `json:"data,omitempty"`
	Error any `json:"error,omitempty"`
}

// Encode renders the data as JSON with a specific status code.
func Encode(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// Wrap in "data" envelope unless it's already an error response
	// (You can simplify this to just strictly wrap everything in "data" if preferred)
	wrapper := Response{Data: data}

	if err := json.NewEncoder(w).Encode(wrapper); err != nil {
		return fmt.Errorf("failed to encode json: %w", err)
	}
	return nil
}

// EncodeError sends a standardized error response
func EncodeError(w http.ResponseWriter, status int, message string, fields map[string]string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	errResp := map[string]any{
		"message": message,
	}
	if len(fields) > 0 {
		errResp["fields"] = fields
	}

	wrapper := Response{Error: errResp}

	if err := json.NewEncoder(w).Encode(wrapper); err != nil {
		return fmt.Errorf("failed to encode json error: %w", err)
	}
	return nil
}

// Decode reads the body and decodes it into the destination struct.
func Decode(w http.ResponseWriter, r *http.Request, dst any) error {
	// Limit body size to 1MB to prevent DOS
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() // Strict contract

	if err := dec.Decode(dst); err != nil {
		return err
	}
	return nil
}
