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

type userRepo struct {
	db *pgxpool.Pool
}

// NewUserRepository creates a new PostgreSQL-backed UserRepository.
func NewUserRepository(db *pgxpool.Pool) repository.UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, full_name, role, email_verified)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.Role,
		user.EmailVerified,
	).Scan(&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if isUniqueViolation(err) {
			return repository.ErrDuplicateEmail
		}
		return fmt.Errorf("userRepo.Create: %w", err)
	}
	return nil
}

func (r *userRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, role, email_verified, created_at, updated_at
		FROM users WHERE id = $1`

	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FullName,
		&user.Role, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("userRepo.GetByID: %w", err)
	}
	return user, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, role, email_verified, created_at, updated_at
		FROM users WHERE email = $1`

	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FullName,
		&user.Role, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("userRepo.GetByEmail: %w", err)
	}
	return user, nil
}

func (r *userRepo) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users SET full_name = $2, email_verified = $3
		WHERE id = $1
		RETURNING updated_at`

	err := r.db.QueryRow(ctx, query, user.ID, user.FullName, user.EmailVerified).
		Scan(&user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.ErrNotFound
		}
		return fmt.Errorf("userRepo.Update: %w", err)
	}
	return nil
}
