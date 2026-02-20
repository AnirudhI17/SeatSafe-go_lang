package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"ticketing/backend/internal/domain"
	"ticketing/backend/internal/dto"
	"ticketing/backend/internal/repository"
	"ticketing/backend/internal/service"
)

// EventHandler handles HTTP requests for event resources.
type EventHandler struct {
	svc *service.EventService
}

func NewEventHandler(svc *service.EventService) *EventHandler {
	return &EventHandler{svc: svc}
}

// CreateEvent godoc
// POST /api/v1/events
func (h *EventHandler) CreateEvent(c *gin.Context) {
	organizerID := mustUserID(c)

	var req dto.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Err(err.Error()))
		return
	}

	// Validate that at least one set of time fields is provided
	if req.StartsAt.IsZero() && req.StartTime.IsZero() {
		c.JSON(http.StatusBadRequest, dto.Err("start time is required (use starts_at or start_time)"))
		return
	}
	if req.EndsAt.IsZero() && req.EndTime.IsZero() {
		c.JSON(http.StatusBadRequest, dto.Err("end time is required (use ends_at or end_time)"))
		return
	}

	event, err := h.svc.CreateEvent(c.Request.Context(), organizerID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Err("failed to create event"))
		return
	}
	c.JSON(http.StatusCreated, dto.OK(event))
}

// GetEvent godoc
// GET /api/v1/events/:id
func (h *EventHandler) GetEvent(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Err("invalid event id"))
		return
	}

	event, err := h.svc.GetEvent(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, dto.Err("event not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.Err("failed to retrieve event"))
		return
	}
	c.JSON(http.StatusOK, dto.OK(event))
}

// ListEvents godoc
// GET /api/v1/events
func (h *EventHandler) ListEvents(c *gin.Context) {
	published := "published"
	statusStr := published
	if s := c.Query("status"); s != "" {
		statusStr = s
	}

	filter := repository.EventFilter{
		Search: c.Query("q"),
	}
	
	// Convert string to domain.EventStatus
	status := domain.EventStatus(statusStr)
	filter.Status = &status

	events, err := h.svc.ListEvents(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Err("failed to list events"))
		return
	}
	c.JSON(http.StatusOK, dto.OK(events))
}

// PublishEvent godoc
// PATCH /api/v1/events/:id/publish
func (h *EventHandler) PublishEvent(c *gin.Context) {
	organizerID := mustUserID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Err("invalid event id"))
		return
	}

	event, err := h.svc.PublishEvent(c.Request.Context(), id, organizerID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, dto.Err("event not found"))
		case errors.Is(err, repository.ErrUnauthorised):
			c.JSON(http.StatusForbidden, dto.Err("not your event"))
		default:
			c.JSON(http.StatusInternalServerError, dto.Err("failed to publish event"))
		}
		return
	}
	c.JSON(http.StatusOK, dto.OK(event))
}
