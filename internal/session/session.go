package session

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"signflow-2026/internal/db"
)

const CookieName = "signflow_session"

const Lifetime = 30 * 24 * time.Hour

var ErrNoSession = errors.New("no session")

type Manager struct {
	q      *db.Queries
	secure bool
}

func NewManager(q *db.Queries, secure bool) *Manager {
	return &Manager{q: q, secure: secure}
}

func (m *Manager) Create(ctx context.Context, w http.ResponseWriter, userID pgtype.UUID) error {
	token, err := randomToken()
	if err != nil {
		return err
	}

	expires := time.Now().Add(Lifetime)

	_, err = m.q.CreateSession(ctx, db.CreateSessionParams{
		TokenHash: hashToken(token),
		UserID:    userID,
		ExpiresAt: pgtype.Timestamptz{Time: expires, Valid: true},
	})
	if err != nil {
		return err
	}

	http.SetCookie(w, m.cookie(token, expires))
	return nil
}

func (m *Manager) User(ctx context.Context, r *http.Request) (db.User, error) {
	c, err := r.Cookie(CookieName)
	if err != nil || c.Value == "" {
		return db.User{}, ErrNoSession
	}

	row, err := m.q.GetUserBySessionToken(ctx, hashToken(c.Value))
	if errors.Is(err, pgx.ErrNoRows) {
		return db.User{}, ErrNoSession
	}
	if err != nil {
		return db.User{}, err
	}

	return row.User, nil
}

func (m *Manager) Destroy(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if c, err := r.Cookie(CookieName); err == nil && c.Value != "" {
		if err := m.q.DeleteSession(ctx, hashToken(c.Value)); err != nil {
			return err
		}
	}

	http.SetCookie(w, m.cookie("", time.Unix(0, 0)))
	return nil
}

func (m *Manager) cookie(value string, expires time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     CookieName,
		Value:    value,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
		Secure:   m.secure,
		SameSite: http.SameSiteLaxMode,
	}
}

func randomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
