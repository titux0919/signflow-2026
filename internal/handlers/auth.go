package handlers

import (
	"errors"
	"net/http"
	"net/mail"
	"strings"

	"github.com/jackc/pgx/v5"

	"signflow-2026/internal/auth"
	"signflow-2026/internal/db"
	"signflow-2026/internal/web"
)

func (h *Handlers) RegisterForm(w http.ResponseWriter, r *http.Request) {
	render(w, r, http.StatusOK, web.Register("", ""))
}

func (h *Handlers) Register(w http.ResponseWriter, r *http.Request) {
	emailAddr := normalizeEmail(r.FormValue("email"))
	password := r.FormValue("password")
	confirm := r.FormValue("confirm")

	fail := func(msg string) {
		render(w, r, http.StatusUnprocessableEntity, web.Register(emailAddr, msg))
	}

	if !validEmail(emailAddr) {
		fail("Please enter a valid email address.")
		return
	}
	if password != confirm {
		fail("Passwords do not match.")
		return
	}

	hash, err := auth.HashPassword(password)
	if errors.Is(err, auth.ErrPasswordLength) {
		fail("Password must be between 8 and 64 characters.")
		return
	}
	if err != nil {
		h.serverError(w, r, "hash password", err)
		return
	}

	// Check for an existing email to show a friendly message.
	if _, err := h.Queries.GetUserByEmail(r.Context(), emailAddr); err == nil {
		fail("That email is already registered. Try logging in.")
		return
	} else if !errors.Is(err, pgx.ErrNoRows) {
		h.serverError(w, r, "lookup user", err)
		return
	}

	user, err := h.Queries.CreateUser(r.Context(), db.CreateUserParams{
		Email:        emailAddr,
		PasswordHash: hash,
	})
	if err != nil {
		if isUniqueViolation(err) {
			fail("That email is already registered. Try logging in.")
			return
		}
		h.serverError(w, r, "create user", err)
		return
	}

	if err := h.Sessions.Create(r.Context(), w, user.ID); err != nil {
		h.serverError(w, r, "create session", err)
		return
	}

	h.Log.Info("user registered", "email", user.Email)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func normalizeEmail(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func validEmail(s string) bool {
	if s == "" || len(s) > 254 {
		return false
	}
	_, err := mail.ParseAddress(s)
	return err == nil
}
