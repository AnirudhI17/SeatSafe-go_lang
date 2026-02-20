package domain

import (
	"time"

	"github.com/google/uuid"
)

// Ticket represents a proof-of-registration for a single seat.
type Ticket struct {
	ID             uuid.UUID  `json:"id"`
	RegistrationID uuid.UUID  `json:"registration_id"`
	EventID        uuid.UUID  `json:"event_id"`
	UserID         uuid.UUID  `json:"user_id"`
	TicketCode     string     `json:"ticket_code"` // e.g. TKT-A3F9-2KXP
	SeatNumber     string     `json:"seat_number,omitempty"`
	IsCheckedIn    bool       `json:"is_checked_in"`
	CheckedInAt    *time.Time `json:"checked_in_at,omitempty"`
	IssuedAt       time.Time  `json:"issued_at"`
	ExpiresAt      *time.Time `json:"expires_at,omitempty"`
}
