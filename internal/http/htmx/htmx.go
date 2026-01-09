package htmx

import (
	"net/http"
)

// IsHTMX checks if the request was initiated by HTMX
func IsHTMX(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

// ReplyRedirect handles client-side redirects via HTMX.
// If you use http.Redirect, HTMX will just swap the content of the button
// with the content of the new page (which looks broken).
// This header tells HTMX to change the browser URL and reload.
func ReplyRedirect(w http.ResponseWriter, path string) {
	w.Header().Set("HX-Redirect", path)
	w.WriteHeader(http.StatusOK) // HTMX expects 200, not 302, for this header to work
}

// ReplyPushURL tells HTMX to update the browser URL bar without a full reload
func ReplyPushURL(w http.ResponseWriter, path string) {
	w.Header().Set("HX-Push-Url", path)
}
