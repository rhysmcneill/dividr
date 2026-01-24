package handler

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/rhysmcneill/dividr/internal/database"
	"github.com/rhysmcneill/dividr/internal/hmrc"
	"github.com/rhysmcneill/dividr/internal/http/handler/render"
	hmrcTempl "github.com/rhysmcneill/dividr/web/templates/hmrc"
)

func (h *Handler) handleAuthDetailsForm(w http.ResponseWriter, r *http.Request) {
	// Render the Templ component
	// Pass "" as error string initially
	_ = render.Component(w, r, http.StatusOK, hmrcTempl.ConnectPage(""))
}

// -- NEW: The Form Submission --
func (h *Handler) handleAuthDetailsSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	nino := r.FormValue("nino")
	mtdID := r.FormValue("mtd_id")

	// 1. Validation (Basic Regex)
	// NINO: 2 letters, 6 numbers, 1 letter (Standard UK format)
	ninoRegex := regexp.MustCompile(`^[A-Za-z]{2}[0-9]{6}[A-Za-z]$`)

	// MTD ID: Usually starts with X, followed by 14 chars total (approx)
	// We'll keep it loose for now to avoid blocking valid IDs, just check length
	if !ninoRegex.MatchString(nino) {
		_ = render.Component(w, r, http.StatusBadRequest, hmrcTempl.ConnectPage("Invalid National Insurance Number format"))
		return
	}
	if len(mtdID) < 10 {
		_ = render.Component(w, r, http.StatusBadRequest, hmrcTempl.ConnectPage("MTD ID seems too short"))
		return
	}

	// 2. Save to Session (Temporary storage)
	h.SessionManager.Put(r.Context(), "temp_nino", nino)
	h.SessionManager.Put(r.Context(), "temp_mtd_id", mtdID)

	// 3. Redirect to the real OAuth start
	http.Redirect(w, r, "/auth/hmrc/start", http.StatusSeeOther)
}

func (h *Handler) handleAuthStart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. Pre-flight Check: Did the user fill in the NINO form?
	// We check our session for the temporary data. If missing, redirect to form.
	if !h.SessionManager.Exists(ctx, "temp_nino") {
		http.Redirect(w, r, "/auth/hmrc/details", http.StatusSeeOther)
		return
	}

	// 2. Initialize the HMRC Service
	authService := hmrc.NewAuthService(h.Config)

	// 3. Generate URL & State
	url, state := authService.GenerateAuthURL()

	// 4. Set Secure Cookie (CSRF Protection for the OAuth flow)
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Expires:  time.Now().Add(10 * time.Minute),
		HttpOnly: true,
		Secure:   h.Config.AppEnv == "production",
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	// 5. Log & Redirect
	slog.Info("starting hmrc oauth flow", "redirect_url", url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// handleAuthCallback processes the response from HMRC after the user logs in.
func (h *Handler) handleAuthCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. Validate State Cookie (CSRF Check)
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil {
		slog.Warn("oauth callback failed: missing state cookie", "error", err)
		http.Error(w, "Session expired, please try again", http.StatusBadRequest)
		return
	}

	queryState := r.URL.Query().Get("state")
	if queryState != stateCookie.Value {
		slog.Warn("oauth callback failed: state mismatch", "query", queryState, "cookie", stateCookie.Value)
		http.Error(w, "Security check failed", http.StatusBadRequest)
		return
	}

	// 2. Get the Authorization Code
	code := r.URL.Query().Get("code")
	if code == "" {
		slog.Warn("oauth callback failed: no code received")
		http.Error(w, "No authorization code received", http.StatusBadRequest)
		return
	}

	// 3. Exchange Code for Access Token
	authService := hmrc.NewAuthService(h.Config)
	token, err := authService.ExchangeCode(ctx, code)
	if err != nil {
		slog.Error("oauth exchange failed", "error", err)
		http.Error(w, "Failed to connect to HMRC", http.StatusInternalServerError)
		return
	}

	// 4. Check User Session (User must be logged in to Dividr)
	// The middleware ensures this usually, but we double-check because we need the UUID.
	userIDStr := h.SessionManager.GetString(ctx, "userID")
	if userIDStr == "" {
		slog.Warn("hmrc callback received for unauthenticated user")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		slog.Error("invalid user uuid in session", "error", err)
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	// 5. Retrieve User Inputs from Session (The "Manual" Lookup)
	userNino := h.SessionManager.GetString(ctx, "temp_nino")
	mtdID := h.SessionManager.GetString(ctx, "temp_mtd_id")

	// If these are missing, the user probably waited too long or opened a new tab.
	if userNino == "" || mtdID == "" {
		slog.Warn("callback session missing nino/mtd_id", "user_id", userID)
		// Send them back to the input form to try again
		http.Redirect(w, r, "/auth/hmrc/details", http.StatusSeeOther)
		return
	}

	// 6. Save to DB
	_, err = h.DB.UpsertHMRCConnection(ctx, database.UpsertHMRCConnectionParams{
		UserID:       database.UUIDToPgtype(userID),
		MtdID:        mtdID,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenExpiry:  database.TimeToPgtype(token.Expiry),
	})

	if err != nil {
		slog.Error("failed to save hmrc connection", "error", err)
		http.Error(w, "Database error saving connection", http.StatusInternalServerError)
		return
	}

	// 7. Cleanup & Logging
	// Remove the temporary inputs now that we have safely stored them in the DB
	h.SessionManager.Remove(ctx, "temp_nino")
	h.SessionManager.Remove(ctx, "temp_mtd_id")

	// Clear the oauth state cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		Path:     "/",
		HttpOnly: true,
		Secure:   h.Config.AppEnv == "production",
	})

	slog.Info("hmrc connected successfully",
		"user_id", userID,
		"mtd_id", mtdID,
	)

	// 8. Redirect to Dashboard
	http.Redirect(w, r, "/app/dashboard", http.StatusSeeOther)
}

// GetValidHMRCToken retrieves a valid access token for a user.
// If the current token is expired, it automatically refreshes it and updates the DB.
func (h *Handler) GetValidHMRCToken(ctx context.Context, userID uuid.UUID) (string, error) {
	// 1. Get current connection details from DB
	conn, err := h.DB.GetHMRCConnectionByUserID(ctx, database.UUIDToPgtype(userID))
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("user not connected to HMRC")
		}
		return "", fmt.Errorf("database error: %w", err)
	}

	// 2. Check Expiry
	// We add a 5-minute buffer. If it expires in < 5 mins, treat it as expired now.
	// This prevents race conditions where the token expires *during* the API call.
	expiry := conn.TokenExpiry.Time
	if time.Now().Add(5 * time.Minute).Before(expiry) {
		// Token is still valid!
		return conn.AccessToken, nil
	}

	slog.Info("hmrc token expired (or close to expiry), refreshing...", "user_id", userID)

	// 3. Token is Expired: Refresh it
	authService := hmrc.NewAuthService(h.Config)

	// We need the refresh token. If it's missing, the user must re-login.
	if conn.RefreshToken == "" {
		slog.Warn("hmrc refresh token missing, user needs to re-authenticate", "user_id", userID)
		return "", fmt.Errorf("refresh token missing, user needs to re-authenticate")
	}

	newToken, err := authService.RefreshAccessToken(conn.RefreshToken)
	if err != nil {
		// If refresh fails (e.g. user revoked access), we might want to log it specifically
		slog.Error("failed to refresh hmrc token", "user_id", userID, "error", err)
		return "", fmt.Errorf("failed to refresh token: %w", err)
	}

	// 4. Update Database
	// We use Upsert to overwrite the old keys with the new ones
	_, err = h.DB.UpsertHMRCConnection(ctx, database.UpsertHMRCConnectionParams{
		UserID:       database.UUIDToPgtype(userID),
		MtdID:        conn.MtdID, // Keep existing MTD ID
		AccessToken:  newToken.AccessToken,
		RefreshToken: newToken.RefreshToken, // HMRC rotates refresh tokens too!
		TokenExpiry:  database.TimeToPgtype(newToken.Expiry),
	})

	if err != nil {
		slog.Error("failed to save refreshed hmrc tokens", "user_id", userID, "error", err)
		return "", fmt.Errorf("failed to save fresh tokens: %w", err)
	}

	slog.Info("hmrc token refreshed successfully", "user_id", userID)

	// 5. Return the new valid Access Token
	return newToken.AccessToken, nil
}
