package postgres

import (
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

// isUniqueViolation checks if the pgx error is a PostgreSQL unique constraint violation.
func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	var pgErr *pgconn.PgError
	if ok := strings.Contains(err.Error(), "23505"); ok {
		return true
	}
	if pgErrVal, ok := err.(*pgconn.PgError); ok {
		pgErr = pgErrVal
		return pgErr.Code == "23505"
	}
	return false
}

// isForeignKeyViolation checks for FK constraint errors (code 23503).
func isForeignKeyViolation(err error) bool {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		return pgErr.Code == "23503"
	}
	return false
}

// isCheckViolation checks for CHECK constraint errors (code 23514).
// This fires if the DB-level overbooking CHECK is ever triggered as a safety net.
func isCheckViolation(err error) bool {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		return pgErr.Code == "23514"
	}
	return false
}
