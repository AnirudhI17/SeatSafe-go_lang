package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"ticketing/backend/internal/domain"
	"ticketing/backend/internal/dto"
	"ticketing/backend/internal/repository"
)

// RegistrationService orchestrates concurrency-safe event booking.
type RegistrationService struct {
	regRepo    repository.RegistrationRepository
	ticketRepo repository.TicketRepository
}

func NewRegistrationService(reg repository.RegistrationRepository, ticket repository.TicketRepository) *RegistrationService {
	return &RegistrationService{regRepo: reg, ticketRepo: ticket}
}

// BookEvent is the service-level booking entry point.
//
// Concurrency Strategy:
//   - Delegates to repository.BookSeat which uses SELECT FOR UPDATE
//   - Retries up to 3 times ONLY for deadlock errors (PG code 40P01)
//   - Business rule errors (ErrEventFull, ErrAlreadyRegistered) fail immediately
//   - Tickets are issued after a successful booking
func (s *RegistrationService) BookEvent(ctx context.Context, eventID, userID uuid.UUID, req dto.BookEventRequest) (*domain.Registration, *domain.Ticket, error) {
	quantity := req.Quantity
	if quantity <= 0 {
		quantity = 1
	}

	const maxRetries = 3
	var (
		reg    *domain.Registration
		err    error
	)

	for attempt := 0; attempt < maxRetries; attempt++ {
		reg, err = s.regRepo.BookSeat(ctx, eventID, userID, quantity)
		if err == nil {
			break
		}

		// Only retry on deadlock — all other errors are non-retryable
		if !isDeadlock(err) {
			return nil, nil, err
		}

		// Exponential backoff: 0ms → 50ms → 100ms
		backoff := time.Duration(attempt*50) * time.Millisecond
		select {
		case <-time.After(backoff):
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		}
	}
	if err != nil {
		return nil, nil, fmt.Errorf("booking failed after %d attempts: %w", maxRetries, err)
	}

	// Issue tickets after successful booking (one ticket per seat)
	tickets := make([]*domain.Ticket, quantity)
	for i := 0; i < quantity; i++ {
		ticket := &domain.Ticket{
			ID:             uuid.New(),
			RegistrationID: reg.ID,
			EventID:        reg.EventID,
			UserID:         reg.UserID,
			TicketCode:     generateTicketCode(),
		}
		if createErr := s.ticketRepo.Create(ctx, ticket); createErr != nil {
			// Non-fatal: registration succeeded. Ticket can be generated async.
			fmt.Printf("WARNING: ticket generation failed for registration %s: %v\n", reg.ID, createErr)
		} else {
			tickets[i] = ticket
		}
	}

	// Return the first ticket for backward compatibility
	var firstTicket *domain.Ticket
	if len(tickets) > 0 && tickets[0] != nil {
		firstTicket = tickets[0]
	}

	return reg, firstTicket, nil
}

func (s *RegistrationService) GetMyRegistrations(ctx context.Context, userID uuid.UUID) ([]*domain.Registration, error) {
	return s.regRepo.ListByUser(ctx, userID)
}

func (s *RegistrationService) CancelRegistration(ctx context.Context, regID, userID uuid.UUID) error {
	return s.regRepo.Cancel(ctx, regID, userID)
}

func (s *RegistrationService) GetMyTickets(ctx context.Context, userID uuid.UUID) ([]*domain.Ticket, error) {
	return s.ticketRepo.ListByUser(ctx, userID)
}

// ListRegistrationsForEvent returns all registrations for a given event.
// Caller is responsible for authorisation (organiser/admin checks).
func (s *RegistrationService) ListRegistrationsForEvent(ctx context.Context, eventID uuid.UUID) ([]*domain.Registration, error) {
	return s.regRepo.ListByEvent(ctx, eventID)
}

// isDeadlock returns true if the error is a PostgreSQL deadlock (code 40P01).
// Deadlocks are transient and safe to retry.
func isDeadlock(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "40P01"
	}
	return false
}
