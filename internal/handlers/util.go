package handlers

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// isUniqueViolation reports whether err is a Postgres unique-constraint error
// (SQLSTATE 23505), so a duplicate-email race becomes a friendly message.
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
