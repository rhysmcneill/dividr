package render

import (
	"net/http"

	"github.com/a-h/templ"
)

// Component renders a Templ component to the response writer.
// It sets the status code and Content-Type header automatically.
func Component(w http.ResponseWriter, r *http.Request, status int, c templ.Component) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)

	// Render the component with the request context
	return c.Render(r.Context(), w)
}
