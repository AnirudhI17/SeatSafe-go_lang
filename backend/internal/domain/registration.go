package domain

import (
	"time"

	"github.com/google/uuid"
)

// RegistrationStatus matches the DB enum lifecycle.
type RegistrationStatus string

const (
	RegistrationStatusPending    RegistrationStatus = "pending"
	RegistrationStatusConfirmed  RegistrationStatus = "confirmed"
	RegistrationStatusCancelled  RegistrationStatus = "cancelled"
	RegistrationStatusWaitlisted RegistrationStatus = "waitlisted"
)

// Registration represents a user's booking for an event.
type Registration struct {
	ID           uuid.UUID          `json:"id"`
	EventID      uuid.UUID          `json:"event_id"`
	UserID       uuid.UUID          `json:"user_id"`
	Status       RegistrationStatus `json:"status"`
	Quantity     int                `json:"quantity"`
	Notes        string             `json:"notes,omitempty"`
	RegisteredAt time.Time          `json:"registered_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
}
