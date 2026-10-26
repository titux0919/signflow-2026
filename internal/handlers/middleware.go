package handlers

import (
	"context"
	"errors"
	"net/http"

	"signflow-2026/internal/db"
	"signflow-2026/internal/session"
)

type ctxKey int

const userCtxKey ctxKey = iota

// LoadUser resolves the session cookie into the current user.
func (h *Handlers) LoadUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := h.Sessions.User(r.Context(), r)
		if err == nil {
			r = r.WithContext(context.WithValue(r.Context(), userCtxKey, user))
		} else if !errors.Is(err, session.ErrNoSession) {
			h.Log.Error("load user from session", "err", err)
		}

		next.ServeHTTP(w, r)
	})
}

// RequireAuth redirects anonymous users to login.
func (h *Handlers) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := UserFrom(r.Context()); !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// UserFrom returns the current user from context.
func UserFrom(ctx context.Context) (db.User, bool) {
	u, ok := ctx.Value(userCtxKey).(db.User)
	return u, ok
}
