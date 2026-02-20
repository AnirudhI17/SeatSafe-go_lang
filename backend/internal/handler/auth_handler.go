package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"ticketing/backend/internal/dto"
	"ticketing/backend/internal/repository"
	"ticketing/backend/internal/service"
)

// AuthHandler handles user registration, login, and profile.
type AuthHandler struct {
	svc *service.UserService
}

func NewAuthHandler(svc *service.UserService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Register godoc
// POST /api/v1/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Err(err.Error()))
		return
	}

	resp, err := h.svc.Register(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateEmail) {
			c.JSON(http.StatusConflict, dto.Err("email address is already registered"))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.Err("registration failed"))
		return
	}

	c.JSON(http.StatusCreated, dto.OK(resp))
}

// Login godoc
// POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Err(err.Error()))
		return
	}

	resp, err := h.svc.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			// Generic message — never reveal whether email or password was wrong
			c.JSON(http.StatusUnauthorized, dto.Err("invalid email or password"))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.Err("login failed"))
		return
	}

	c.JSON(http.StatusOK, dto.OK(resp))
}

// Profile godoc
// GET /api/v1/auth/me  (protected)
func (h *AuthHandler) Profile(c *gin.Context) {
	rawID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.Err("not authenticated"))
		return
	}
	userID, ok := rawID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.Err("invalid session"))
		return
	}

	user, err := h.svc.GetProfile(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, dto.Err("user not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.Err("failed to fetch profile"))
		return
	}

	c.JSON(http.StatusOK, dto.OK(dto.UserProfile{
		ID:       user.ID.String(),
		Email:    user.Email,
		FullName: user.FullName,
		Role:     string(user.Role),
	}))
}
