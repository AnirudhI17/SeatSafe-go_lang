package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"ticketing/backend/internal/dto"
	"ticketing/backend/internal/repository"
)

// ErrorHandler is a centralised error-mapping middleware.
// It wraps every handler response: if a handler sets a "error" key in the
// context (via c.Error(err)), this middleware converts it to a proper HTTP response.
// Direct c.JSON(...) in handlers is still supported alongside this.
func ErrorHandler(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Only process if there are errors AND no response was written yet
		if len(c.Errors) == 0 || c.Writer.Written() {
			return
		}

		err := c.Errors.Last().Err
		status, message := classifyError(err, logger)
		c.JSON(status, dto.Err(message))
	}
}

// classifyError maps domain/infrastructure errors to HTTP status + message.
func classifyError(err error, logger *zap.Logger) (int, string) {
	// ── Domain errors ────────────────────────────────────────────────────────
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return http.StatusNotFound, "resource not found"
	case errors.Is(err, repository.ErrEventFull):
		return http.StatusConflict, "event is fully booked"
	case errors.Is(err, repository.ErrAlreadyRegistered):
		return http.StatusConflict, "you are already registered for this event"
	case errors.Is(err, repository.ErrEventNotPublished):
		return http.StatusUnprocessableEntity, "event is not open for registration"
	case errors.Is(err, repository.ErrUnauthorised):
		return http.StatusForbidden, "not authorised to perform this action"
	case errors.Is(err, repository.ErrDuplicateEmail):
		return http.StatusConflict, "email address is already registered"
	}

	// ── PostgreSQL errors ─────────────────────────────────────────────────────
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "40P01": // deadlock_detected
			logger.Warn("deadlock detected", zap.Error(err))
			return http.StatusServiceUnavailable, "please try again"
		case "23514": // check_violation — overbooking safety net
			logger.Error("CHECK constraint violated — overbooking safety net fired", zap.Error(err))
			return http.StatusConflict, "booking capacity exceeded"
		default:
			logger.Error("unhandled postgres error", zap.String("code", pgErr.Code), zap.Error(err))
			return http.StatusInternalServerError, "database error"
		}
	}

	// ── Catch-all ─────────────────────────────────────────────────────────────
	logger.Error("unhandled error", zap.Error(err))
	return http.StatusInternalServerError, "an internal error occurred"
}
