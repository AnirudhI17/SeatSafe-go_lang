package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"ticketing/backend/internal/dto"
)

// AuthMiddleware validates the JWT Bearer token in the Authorization header.
// On success, sets "userID" (uuid.UUID) and "userRole" (string) in the Gin context.
// Full JWT implementation is in Phase 6; this stub enables protected routes in Phase 3.
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Err("missing or invalid authorization header"))
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Err("invalid or expired token"))
			return
		}

		rawID, _ := claims["sub"].(string)
		userID, err := uuid.Parse(rawID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Err("invalid token subject"))
			return
		}

		c.Set("userID", userID)
		c.Set("userRole", claims["role"])
		c.Next()
	}
}
