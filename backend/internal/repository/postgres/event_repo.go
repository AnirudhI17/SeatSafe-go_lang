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

type eventRepo struct {
	db *pgxpool.Pool
}

// NewEventRepository creates a new PostgreSQL-backed EventRepository.
func NewEventRepository(db *pgxpool.Pool) repository.EventRepository {
	return &eventRepo{db: db}
}

func (r *eventRepo) Create(ctx context.Context, event *domain.Event) error {
	query := `
		INSERT INTO events
			(id, organizer_id, title, description, location, is_online, online_url,
			 starts_at, ends_at, capacity, price_cents, currency, banner_url, status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
		RETURNING registered_count, created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		event.ID, event.OrganizerID, event.Title, event.Description,
		event.Location, event.IsOnline, event.OnlineURL,
		event.StartsAt, event.EndsAt, event.Capacity,
		event.PriceCents, event.Currency, event.BannerURL, event.Status,
	).Scan(&event.RegisteredCount, &event.CreatedAt, &event.UpdatedAt)
}

func (r *eventRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	query := `
		SELECT id, organizer_id, title, COALESCE(description, ''), location, is_online, COALESCE(online_url, ''),
		       starts_at, ends_at, capacity, registered_count, price_cents, currency,
		       COALESCE(banner_url, ''), status, created_at, updated_at
		FROM events WHERE id = $1`

	e := &domain.Event{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&e.ID, &e.OrganizerID, &e.Title, &e.Description, &e.Location,
		&e.IsOnline, &e.OnlineURL, &e.StartsAt, &e.EndsAt, &e.Capacity,
		&e.RegisteredCount, &e.PriceCents, &e.Currency, &e.BannerURL,
		&e.Status, &e.CreatedAt, &e.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("eventRepo.GetByID: %w", err)
	}
	return e, nil
}

func (r *eventRepo) List(ctx context.Context, filter repository.EventFilter) ([]*domain.Event, error) {
	clauses := []string{"1=1"}
	args := []any{}
	argIdx := 1

	if filter.Status != nil {
		clauses = append(clauses, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, *filter.Status)
		argIdx++
	}

	if filter.Search != "" {
		clauses = append(clauses, fmt.Sprintf(
			"to_tsvector('english', title || ' ' || COALESCE(description, '')) @@ plainto_tsquery('english', $%d)", argIdx))
		args = append(args, filter.Search)
		argIdx++
	}

	if filter.PageSize == 0 {
		filter.PageSize = 20
	}
	offset := filter.Page * filter.PageSize

	args = append(args, filter.PageSize, offset)
	query := fmt.Sprintf(`
		SELECT id, organizer_id, title, COALESCE(description, ''), location, is_online, COALESCE(online_url, ''),
		       starts_at, ends_at, capacity, registered_count, price_cents, currency,
		       COALESCE(banner_url, ''), status, created_at, updated_at
		FROM events
		WHERE %s
		ORDER BY starts_at ASC
		LIMIT $%d OFFSET $%d`,
		strings.Join(clauses, " AND "), argIdx, argIdx+1)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("eventRepo.List: %w", err)
	}
	defer rows.Close()

	var events []*domain.Event
	for rows.Next() {
		e := &domain.Event{}
		if err := rows.Scan(
			&e.ID, &e.OrganizerID, &e.Title, &e.Description, &e.Location,
			&e.IsOnline, &e.OnlineURL, &e.StartsAt, &e.EndsAt, &e.Capacity,
			&e.RegisteredCount, &e.PriceCents, &e.Currency, &e.BannerURL,
			&e.Status, &e.CreatedAt, &e.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("eventRepo.List scan: %w", err)
		}
		events = append(events, e)
	}
	return events, rows.Err()
}

func (r *eventRepo) Update(ctx context.Context, event *domain.Event) error {
	query := `
		UPDATE events
		SET title=$2, description=$3, location=$4, is_online=$5, online_url=$6,
		    starts_at=$7, ends_at=$8, capacity=$9, banner_url=$10, status=$11
		WHERE id=$1
		RETURNING updated_at`

	err := r.db.QueryRow(ctx, query,
		event.ID, event.Title, event.Description, event.Location,
		event.IsOnline, event.OnlineURL, event.StartsAt, event.EndsAt,
		event.Capacity, event.BannerURL, event.Status,
	).Scan(&event.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.ErrNotFound
		}
		return fmt.Errorf("eventRepo.Update: %w", err)
	}
	return nil
}

func (r *eventRepo) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM events WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("eventRepo.Delete: %w", err)
	}
	if result.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}
