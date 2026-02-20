package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"ticketing/backend/internal/domain"
	"ticketing/backend/internal/repository"
	"ticketing/backend/internal/dto"
)

// EventService holds business logic for events.
type EventService struct {
	eventRepo repository.EventRepository
}

func NewEventService(repo repository.EventRepository) *EventService {
	return &EventService{eventRepo: repo}
}

func (s *EventService) CreateEvent(ctx context.Context, organizerID uuid.UUID, req dto.CreateEventRequest) (*domain.Event, error) {
	// Handle alternative field names for backward compatibility
	startsAt := req.StartsAt
	endsAt := req.EndsAt
	priceCents := req.PriceCents
	
	// If alternative fields are provided, use them
	if !req.StartTime.IsZero() {
		startsAt = req.StartTime
	}
	if !req.EndTime.IsZero() {
		endsAt = req.EndTime
	}
	if req.Price > 0 {
		priceCents = int(req.Price * 100) // Convert dollars to cents
	}
	
	event := &domain.Event{
		ID:          uuid.New(),
		OrganizerID: organizerID,
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		IsOnline:    req.IsOnline,
		OnlineURL:   req.OnlineURL,
		StartsAt:    startsAt,
		EndsAt:      endsAt,
		Capacity:    req.Capacity,
		PriceCents:  priceCents,
		Currency:    "USD",
		BannerURL:   req.BannerURL,
		Status:      domain.EventStatusDraft,
	}

	if err := s.eventRepo.Create(ctx, event); err != nil {
		return nil, fmt.Errorf("EventService.CreateEvent: %w", err)
	}
	return event, nil
}

func (s *EventService) GetEvent(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	return s.eventRepo.GetByID(ctx, id)
}

func (s *EventService) ListEvents(ctx context.Context, filter repository.EventFilter) ([]*domain.Event, error) {
	return s.eventRepo.List(ctx, filter)
}

func (s *EventService) PublishEvent(ctx context.Context, id, organizerID uuid.UUID) (*domain.Event, error) {
	event, err := s.eventRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if event.OrganizerID != organizerID {
		return nil, repository.ErrUnauthorised
	}
	event.Status = domain.EventStatusPublished
	if err := s.eventRepo.Update(ctx, event); err != nil {
		return nil, err
	}
	return event, nil
}

func (s *EventService) UpdateEvent(ctx context.Context, id, organizerID uuid.UUID, req dto.UpdateEventRequest) (*domain.Event, error) {
	event, err := s.eventRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if event.OrganizerID != organizerID {
		return nil, repository.ErrUnauthorised
	}
	if req.Title != "" {
		event.Title = req.Title
	}
	if req.Description != "" {
		event.Description = req.Description
	}
	if req.Capacity > 0 {
		event.Capacity = req.Capacity
	}
	if err := s.eventRepo.Update(ctx, event); err != nil {
		return nil, err
	}
	return event, nil
}
