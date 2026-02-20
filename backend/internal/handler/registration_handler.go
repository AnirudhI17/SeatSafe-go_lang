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

// RegistrationHandler handles HTTP requests for booking and tickets.
type RegistrationHandler struct {
	svc *service.RegistrationService
}

func NewRegistrationHandler(svc *service.RegistrationService) *RegistrationHandler {
	return &RegistrationHandler{svc: svc}
}

// BookEvent godoc
// POST /api/v1/events/:id/register
func (h *RegistrationHandler) BookEvent(c *gin.Context) {
	userID := mustUserID(c)
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Err("invalid event id"))
		return
	}

	var req dto.BookEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Err(err.Error()))
		return
	}

	reg, ticket, err := h.svc.BookEvent(c.Request.Context(), eventID, userID, req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrEventFull):
			// 409 Conflict — standard for "resource state prevents the action"
			c.JSON(http.StatusConflict, dto.Err("event is fully booked"))
		case errors.Is(err, repository.ErrAlreadyRegistered):
			c.JSON(http.StatusConflict, dto.Err("you are already registered for this event"))
		case errors.Is(err, repository.ErrEventNotPublished):
			c.JSON(http.StatusUnprocessableEntity, dto.Err("event is not open for registration"))
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, dto.Err("event not found"))
		default:
			c.JSON(http.StatusInternalServerError, dto.Err("booking failed"))
		}
		return
	}

	c.JSON(http.StatusCreated, dto.OK(gin.H{
		"registration": reg,
		"ticket":       ticket,
	}))
}

// MyRegistrations godoc
// GET /api/v1/registrations/me
func (h *RegistrationHandler) MyRegistrations(c *gin.Context) {
	userID := mustUserID(c)
	regs, err := h.svc.GetMyRegistrations(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Err("failed to fetch registrations"))
		return
	}
	c.JSON(http.StatusOK, dto.OK(regs))
}

// MyTickets godoc
// GET /api/v1/tickets/me
func (h *RegistrationHandler) MyTickets(c *gin.Context) {
	userID := mustUserID(c)
	tickets, err := h.svc.GetMyTickets(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Err("failed to fetch tickets"))
		return
	}
	c.JSON(http.StatusOK, dto.OK(tickets))
}

// CancelRegistration godoc
// DELETE /api/v1/registrations/:id
func (h *RegistrationHandler) CancelRegistration(c *gin.Context) {
	userID := mustUserID(c)
	regID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Err("invalid registration id"))
		return
	}

	if err := h.svc.CancelRegistration(c.Request.Context(), regID, userID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, dto.Err("registration not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.Err("cancellation failed"))
		return
	}
	c.JSON(http.StatusOK, dto.APIResponse{Success: true, Message: "registration cancelled"})
}

// ListEventRegistrations godoc
// GET /api/v1/events/:id/registrations  (organiser/admin only)
func (h *RegistrationHandler) ListEventRegistrations(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Err("invalid event id"))
		return
	}

	regs, err := h.svc.ListRegistrationsForEvent(c.Request.Context(), eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Err("failed to fetch registrations"))
		return
	}
	c.JSON(http.StatusOK, dto.OK(regs))
}

// mustUserID is a helper to extract the authenticated user's ID set by auth middleware.
func mustUserID(c *gin.Context) uuid.UUID {
	raw, _ := c.Get("userID")
	if id, ok := raw.(uuid.UUID); ok {
		return id
	}
	return uuid.Nil
}
