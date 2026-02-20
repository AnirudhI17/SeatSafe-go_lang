package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"ticketing/backend/internal/domain"
	"ticketing/backend/internal/repository"
)

type registrationRepo struct {
	db *pgxpool.Pool
}

// NewRegistrationRepository creates a new PostgreSQL-backed RegistrationRepository.
func NewRegistrationRepository(db *pgxpool.Pool) repository.RegistrationRepository {
	return &registrationRepo{db: db}
}

// BookSeat is the concurrency-safe booking method.
//
// Concurrency Strategy: Pessimistic Row-Level Locking (SELECT FOR UPDATE)
//
// When two goroutines attempt to book the last seat simultaneously:
//  1. Both open a transaction (BEGIN)
//  2. Both execute SELECT ... FOR UPDATE on the same event row
//  3. PostgreSQL grants the lock to ONE goroutine; the other BLOCKS and waits
//  4. The first goroutine checks capacity → seats available → inserts + increments
//  5. First goroutine commits; lock is released
//  6. The second goroutine's SELECT FOR UPDATE now returns the UPDATED row
//     (registered_count is already incremented)
//  7. Second goroutine checks capacity → no seats → returns ErrEventFull
//
// This is implemented fully in Phase 4. This stub exists for Phase 3 compilation.
func (r *registrationRepo) BookSeat(ctx context.Context, eventID, userID uuid.UUID, quantity int) (*domain.Registration, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("BookSeat: begin tx: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	// ── Step 1: Lock the event row ────────────────────────────────────────────
	// SELECT FOR UPDATE ensures no two concurrent transactions can read
	// stale capacity values simultaneously.
	var capacity, registeredCount int
	var status domain.EventStatus
	lockQuery := `
		SELECT capacity, registered_count, status
		FROM events
		WHERE id = $1
		FOR UPDATE`

	err = tx.QueryRow(ctx, lockQuery, eventID).Scan(&capacity, &registeredCount, &status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("BookSeat: lock event row: %w", err)
	}

	// ── Step 2: Business rule checks (inside locked transaction) ──────────────
	if status != domain.EventStatusPublished {
		return nil, repository.ErrEventNotPublished
	}
	if registeredCount+quantity > capacity {
		return nil, repository.ErrEventFull
	}

	// ── Step 3: Check for duplicate active registration ───────────────────────
	var existingID uuid.UUID
	dupQuery := `
		SELECT id FROM registrations
		WHERE user_id = $1 AND event_id = $2
		  AND status IN ('pending','confirmed','waitlisted')
		LIMIT 1`
	err = tx.QueryRow(ctx, dupQuery, userID, eventID).Scan(&existingID)
	if err == nil {
		return nil, repository.ErrAlreadyRegistered
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("BookSeat: duplicate check: %w", err)
	}

	// ── Step 4: Insert the registration ──────────────────────────────────────
	reg := &domain.Registration{
		ID:       uuid.New(),
		EventID:  eventID,
		UserID:   userID,
		Status:   domain.RegistrationStatusConfirmed,
		Quantity: quantity,
	}
	insertQuery := `
		INSERT INTO registrations (id, event_id, user_id, status, quantity)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING registered_at, updated_at`

	err = tx.QueryRow(ctx, insertQuery,
		reg.ID, reg.EventID, reg.UserID, reg.Status, reg.Quantity,
	).Scan(&reg.RegisteredAt, &reg.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, repository.ErrAlreadyRegistered
		}
		return nil, fmt.Errorf("BookSeat: insert registration: %w", err)
	}

	// ── Step 5: Atomically increment registered_count ─────────────────────────
	// This is safe because we hold the row lock.
	updateQuery := `
		UPDATE events
		SET registered_count = registered_count + $2
		WHERE id = $1`
	if _, err = tx.Exec(ctx, updateQuery, eventID, quantity); err != nil {
		return nil, fmt.Errorf("BookSeat: update registered_count: %w", err)
	}

	// ── Step 6: Commit — release the row lock ─────────────────────────────────
	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("BookSeat: commit: %w", err)
	}

	return reg, nil
}

func (r *registrationRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Registration, error) {
	query := `
		SELECT id, event_id, user_id, status, quantity, COALESCE(notes, ''), registered_at, updated_at
		FROM registrations WHERE id = $1`

	reg := &domain.Registration{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&reg.ID, &reg.EventID, &reg.UserID, &reg.Status,
		&reg.Quantity, &reg.Notes, &reg.RegisteredAt, &reg.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("registrationRepo.GetByID: %w", err)
	}
	return reg, nil
}

func (r *registrationRepo) ListByUser(ctx context.Context, userID uuid.UUID) ([]*domain.Registration, error) {
	return r.listWith(ctx, "user_id = $1", userID)
}

func (r *registrationRepo) ListByEvent(ctx context.Context, eventID uuid.UUID) ([]*domain.Registration, error) {
	return r.listWith(ctx, "event_id = $1", eventID)
}

func (r *registrationRepo) listWith(ctx context.Context, where string, arg any) ([]*domain.Registration, error) {
	query := fmt.Sprintf(`
		SELECT id, event_id, user_id, status, quantity, COALESCE(notes, ''), registered_at, updated_at
		FROM registrations WHERE %s ORDER BY registered_at DESC`, where)

	rows, err := r.db.Query(ctx, query, arg)
	if err != nil {
		return nil, fmt.Errorf("registrationRepo.listWith: %w", err)
	}
	defer rows.Close()

	var list []*domain.Registration
	for rows.Next() {
		reg := &domain.Registration{}
		if err := rows.Scan(
			&reg.ID, &reg.EventID, &reg.UserID, &reg.Status,
			&reg.Quantity, &reg.Notes, &reg.RegisteredAt, &reg.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, reg)
	}
	return list, rows.Err()
}

func (r *registrationRepo) Cancel(ctx context.Context, id, userID uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("Cancel: begin tx: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	var quantity int
	var eventID uuid.UUID
	var status domain.RegistrationStatus

	checkQuery := `SELECT event_id, quantity, status FROM registrations WHERE id = $1 AND user_id = $2 FOR UPDATE`
	err = tx.QueryRow(ctx, checkQuery, id, userID).Scan(&eventID, &quantity, &status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.ErrNotFound
		}
		return fmt.Errorf("Cancel: lock registration: %w", err)
	}

	if status == domain.RegistrationStatusCancelled {
		return nil // idempotent
	}

	if _, err = tx.Exec(ctx,
		`UPDATE registrations SET status = 'cancelled' WHERE id = $1`, id,
	); err != nil {
		return fmt.Errorf("Cancel: update status: %w", err)
	}

	if _, err = tx.Exec(ctx,
		`UPDATE events SET registered_count = registered_count - $2 WHERE id = $1`, eventID, quantity,
	); err != nil {
		return fmt.Errorf("Cancel: decrement count: %w", err)
	}

	_ = strings.TrimSpace("") // suppress unused import lint
	return tx.Commit(ctx)
}
