package handler

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/rhysmcneill/dividr/internal/database"
	"github.com/rhysmcneill/dividr/internal/http/handler/render"
	"github.com/rhysmcneill/dividr/web/templates/auth"
	"golang.org/x/crypto/bcrypt"
)

// --- RENDER HANDLERS ---

func (h *Handler) handleLoginPage(w http.ResponseWriter, r *http.Request) {
	if err := render.Component(w, r, http.StatusOK, auth.Login()); err != nil {
		h.respondWithError(w, r, err)
		return
	}
}

func (h *Handler) handleSignupPage(w http.ResponseWriter, r *http.Request) {
	if err := render.Component(w, r, http.StatusOK, auth.Signup()); err != nil {
		h.respondWithError(w, r, err)
		return
	}
}

// --- LOGIC HANDLERS ---

func (h *Handler) handleSignupSubmit(w http.ResponseWriter, r *http.Request) {
	// 1. Parse Form
	email := r.FormValue("email")
	password := r.FormValue("password")

	// 2. Hash Password
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		h.respondWithError(w, r, err)
		return
	}

	// 3. Insert into DB (using generated SQLC queries)
	q := database.New(h.DB.Pool)
	user, err := q.CreateUser(r.Context(), database.CreateUserParams{
		Email:        email,
		PasswordHash: string(hashedBytes),
	})
	if err != nil {
		// Check for duplicate email error (SQLSTATE 23505)
		// For now, generic error
		h.respondWithError(w, r, err)
		return
	}

	// 4. Log them in automatically
	// "RenewToken" is a security best practice when privilege changes (prevents fixation attacks)
	if err := h.SessionManager.RenewToken(r.Context()); err != nil {
		h.respondWithError(w, r, err)
		return
	}
	// Store User ID in session
	h.SessionManager.Put(r.Context(), "userID", user.ID.String())

	// 5. Redirect to Dashboard
	http.Redirect(w, r, "/app/dashboard", http.StatusSeeOther)
}

func (h *Handler) handleLoginSubmit(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	// 1. Find User
	q := database.New(h.DB.Pool)
	user, err := q.GetUserByEmail(r.Context(), email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Generic error message for security
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}
		h.respondWithError(w, r, err)
		return
	}

	// 2. Check Password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// 3. Create Session
	if err := h.SessionManager.RenewToken(r.Context()); err != nil {
		h.respondWithError(w, r, err)
		return
	}
	h.SessionManager.Put(r.Context(), "userID", user.ID.String())
	http.Redirect(w, r, "/app/dashboard", http.StatusSeeOther)
}

func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) {
	// 1. Renew the token (security best practice to prevent session fixation)
	_ = h.SessionManager.RenewToken(r.Context())

	// 2. Clean up entire session
	_ = h.SessionManager.Destroy(r.Context())

	// 3. Redirect to Home/Login
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
