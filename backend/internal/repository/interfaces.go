package repository

import (
	"context"

	"github.com/google/uuid"
	"ticketing/backend/internal/domain"
)

// UserRepository defines the contract for user persistence operations.
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
}

// EventRepository defines the contract for event persistence operations.
type EventRepository interface {
	Create(ctx context.Context, event *domain.Event) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error)
	List(ctx context.Context, filter EventFilter) ([]*domain.Event, error)
	Update(ctx context.Context, event *domain.Event) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// RegistrationRepository defines the contract for booking persistence.
// The concurrency-safe booking method uses a DB transaction internally.
type RegistrationRepository interface {
	// BookSeat is the critical concurrency-safe method.
	// It opens a transaction, locks the event row with SELECT FOR UPDATE,
	// checks capacity, inserts the registration, and increments registered_count.
	// Returns ErrEventFull if no seats remain, ErrAlreadyRegistered if duplicate.
	BookSeat(ctx context.Context, eventID, userID uuid.UUID, quantity int) (*domain.Registration, error)

	GetByID(ctx context.Context, id uuid.UUID) (*domain.Registration, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*domain.Registration, error)
	ListByEvent(ctx context.Context, eventID uuid.UUID) ([]*domain.Registration, error)
	Cancel(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

// TicketRepository defines the contract for ticket persistence.
type TicketRepository interface {
	Create(ctx context.Context, ticket *domain.Ticket) error
	GetByCode(ctx context.Context, code string) (*domain.Ticket, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*domain.Ticket, error)
	CheckIn(ctx context.Context, ticketCode string) (*domain.Ticket, error)
}

// EventFilter holds optional query parameters for listing events.
type EventFilter struct {
	Status   *domain.EventStatus
	Search   string
	Page     int
	PageSize int
}
