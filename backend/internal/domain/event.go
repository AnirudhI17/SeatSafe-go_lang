package domain

import (
	"time"

	"github.com/google/uuid"
)

// EventStatus matches the DB enum.
type EventStatus string

const (
	EventStatusDraft     EventStatus = "draft"
	EventStatusPublished EventStatus = "published"
	EventStatusCancelled EventStatus = "cancelled"
	EventStatusCompleted EventStatus = "completed"
)

// Event is the core entity representing a ticketed event.
type Event struct {
	ID              uuid.UUID   `json:"id"`
	OrganizerID     uuid.UUID   `json:"organizer_id"`
	Title           string      `json:"title"`
	Description     string      `json:"description"`
	Location        string      `json:"location"`
	IsOnline        bool        `json:"is_online"`
	OnlineURL       string      `json:"online_url,omitempty"`
	StartsAt        time.Time   `json:"starts_at"`
	EndsAt          time.Time   `json:"ends_at"`
	Capacity        int         `json:"capacity"`
	RegisteredCount int         `json:"registered_count"`
	PriceCents      int         `json:"price_cents"`
	Currency        string      `json:"currency"`
	BannerURL       string      `json:"banner_url,omitempty"`
	Status          EventStatus `json:"status"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

// AvailableSeats returns remaining capacity.
func (e *Event) AvailableSeats() int {
	return e.Capacity - e.RegisteredCount
}

// IsFull returns true if no seats remain.
func (e *Event) IsFull() bool {
	return e.RegisteredCount >= e.Capacity
}

// IsFree returns true if the event is free of charge.
func (e *Event) IsFree() bool {
	return e.PriceCents == 0
}
