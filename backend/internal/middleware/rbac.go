package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ticketing/backend/internal/dto"
)

// RequireRole returns a middleware that enforces role-based access control.
// Must be used after AuthMiddleware (which sets the "userRole" context key).
//
// Usage:
//
//	auth.POST("/events", RequireRole("organizer", "admin"), eventH.CreateEvent)
func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}

	return func(c *gin.Context) {
		role, exists := c.Get("userRole")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Err("authentication required"))
			return
		}
		roleStr, ok := role.(string)
		if !ok || !allowed[roleStr] {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.Err("insufficient permissions"))
			return
		}
		c.Next()
	}
}
