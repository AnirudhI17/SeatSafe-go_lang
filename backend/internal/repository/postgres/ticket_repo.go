package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"ticketing/backend/internal/domain"
	"ticketing/backend/internal/repository"
)

type ticketRepo struct {
	db *pgxpool.Pool
}

// NewTicketRepository creates a new PostgreSQL-backed TicketRepository.
func NewTicketRepository(db *pgxpool.Pool) repository.TicketRepository {
	return &ticketRepo{db: db}
}

func (r *ticketRepo) Create(ctx context.Context, t *domain.Ticket) error {
	query := `
		INSERT INTO tickets (id, registration_id, event_id, user_id, ticket_code, seat_number)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING issued_at`

	return r.db.QueryRow(ctx, query,
		t.ID, t.RegistrationID, t.EventID, t.UserID, t.TicketCode, t.SeatNumber,
	).Scan(&t.IssuedAt)
}

func (r *ticketRepo) GetByCode(ctx context.Context, code string) (*domain.Ticket, error) {
	query := `
		SELECT id, registration_id, event_id, user_id, ticket_code, seat_number,
		       is_checked_in, checked_in_at, issued_at, expires_at
		FROM tickets WHERE ticket_code = $1`

	t := &domain.Ticket{}
	err := r.db.QueryRow(ctx, query, code).Scan(
		&t.ID, &t.RegistrationID, &t.EventID, &t.UserID, &t.TicketCode,
		&t.SeatNumber, &t.IsCheckedIn, &t.CheckedInAt, &t.IssuedAt, &t.ExpiresAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("ticketRepo.GetByCode: %w", err)
	}
	return t, nil
}

func (r *ticketRepo) ListByUser(ctx context.Context, userID uuid.UUID) ([]*domain.Ticket, error) {
	query := `
		SELECT id, registration_id, event_id, user_id, ticket_code, seat_number,
		       is_checked_in, checked_in_at, issued_at, expires_at
		FROM tickets WHERE user_id = $1 ORDER BY issued_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("ticketRepo.ListByUser: %w", err)
	}
	defer rows.Close()

	var list []*domain.Ticket
	for rows.Next() {
		t := &domain.Ticket{}
		if err := rows.Scan(
			&t.ID, &t.RegistrationID, &t.EventID, &t.UserID, &t.TicketCode,
			&t.SeatNumber, &t.IsCheckedIn, &t.CheckedInAt, &t.IssuedAt, &t.ExpiresAt,
		); err != nil {
			return nil, err
		}
		list = append(list, t)
	}
	return list, rows.Err()
}

func (r *ticketRepo) CheckIn(ctx context.Context, ticketCode string) (*domain.Ticket, error) {
	query := `
		UPDATE tickets
		SET is_checked_in = TRUE, checked_in_at = NOW()
		WHERE ticket_code = $1 AND is_checked_in = FALSE
		RETURNING id, registration_id, event_id, user_id, ticket_code, seat_number,
		          is_checked_in, checked_in_at, issued_at, expires_at`

	t := &domain.Ticket{}
	err := r.db.QueryRow(ctx, query, ticketCode).Scan(
		&t.ID, &t.RegistrationID, &t.EventID, &t.UserID, &t.TicketCode,
		&t.SeatNumber, &t.IsCheckedIn, &t.CheckedInAt, &t.IssuedAt, &t.ExpiresAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound // already checked in or doesn't exist
		}
		return nil, fmt.Errorf("ticketRepo.CheckIn: %w", err)
	}
	return t, nil
}
