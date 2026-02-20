//go:build integration

package integration

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"ticketing/backend/internal/config"
	"ticketing/backend/internal/repository/postgres"
)

// setupTestDB creates a pgxpool connection using DATABASE_URL from the environment.
// Tests are expected to run against your Supabase development branch or a local DB.
func setupTestDB(t *testing.T, ctx context.Context) *pgxpool.Pool {
	t.Helper()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Fallback: try loading from .env in backend root
		if url := loadEnvFile("../../.env", "DATABASE_URL"); url != "" {
			dbURL = url
		}
	}
	if dbURL == "" {
		t.Fatal("DATABASE_URL environment variable is required for integration tests")
	}

	cfg := config.DatabaseConfig{
		URL:             dbURL,
		MaxConns:        20, // allow high concurrency
		MinConns:        5,
		MaxConnIdleTime: 5 * time.Minute,
	}

	pool, err := postgres.NewPool(ctx, cfg)
	if err != nil {
		t.Fatalf("setupTestDB: failed to connect: %v", err)
	}
	return pool
}

// cleanupTestDB closes the pool.
func cleanupTestDB(t *testing.T, ctx context.Context, db *pgxpool.Pool) {
	t.Helper()
	db.Close()
}

// seedUser inserts a test user and returns their UUID.
func seedUser(t *testing.T, ctx context.Context, db *pgxpool.Pool, email string) uuid.UUID {
	t.Helper()

	// Use ON CONFLICT to make seeds idempotent when re-running tests.
	id := uuid.New()
	_, err := db.Exec(ctx, `
		INSERT INTO users (id, email, password_hash, full_name, role)
		VALUES ($1, $2, '$2a$12$testhashdoesnotmatter', 'Test User', 'attendee')
		ON CONFLICT (email) DO UPDATE SET id = EXCLUDED.id
		RETURNING id
	`, id, email)

	// Fetch the actual ID (might have been set by ON CONFLICT)
	var actualID uuid.UUID
	err = db.QueryRow(ctx, `SELECT id FROM users WHERE email = $1`, email).Scan(&actualID)
	if err != nil {
		t.Fatalf("seedUser: %v", err)
	}
	return actualID
}

// seedEvent inserts a published test event with the given capacity.
func seedEvent(t *testing.T, ctx context.Context, db *pgxpool.Pool, organizerID uuid.UUID, capacity int) uuid.UUID {
	t.Helper()

	id := uuid.New()
	now := time.Now()
	_, err := db.Exec(ctx, `
		INSERT INTO events
			(id, organizer_id, title, description, location, starts_at, ends_at,
			 capacity, registered_count, price_cents, currency, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 0, 0, 'USD', 'published')
	`,
		id,
		organizerID,
		fmt.Sprintf("Test Event %s", id.String()[:8]),
		"Concurrency test event",
		"Test Venue",
		now.Add(24*time.Hour),
		now.Add(48*time.Hour),
		capacity,
	)
	if err != nil {
		t.Fatalf("seedEvent: %v", err)
	}
	return id
}

// assertRegisteredCount verifies the events.registered_count in the DB.
func assertRegisteredCount(t *testing.T, ctx context.Context, db *pgxpool.Pool, eventID uuid.UUID, expected int) {
	t.Helper()

	var count int
	err := db.QueryRow(ctx, `SELECT registered_count FROM events WHERE id = $1`, eventID).Scan(&count)
	if err != nil {
		t.Fatalf("assertRegisteredCount: query failed: %v", err)
	}
	if count != expected {
		t.Errorf("DB registered_count = %d, want %d — OVERBOOKING DETECTED", count, expected)
	}
}

// assertRegistrationRows verifies the actual number of confirmed registration rows.
func assertRegistrationRows(t *testing.T, ctx context.Context, db *pgxpool.Pool, eventID uuid.UUID, expected int) {
	t.Helper()

	var count int
	err := db.QueryRow(ctx, `
		SELECT COUNT(*) FROM registrations
		WHERE event_id = $1 AND status IN ('confirmed','pending')
	`, eventID).Scan(&count)
	if err != nil {
		t.Fatalf("assertRegistrationRows: query failed: %v", err)
	}
	if count != expected {
		t.Errorf("registration table has %d rows for event, want %d", count, expected)
	}
}

// loadEnvFile is a minimal .env reader used as fallback in test setup.
func loadEnvFile(path, key string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	for _, line := range splitLines(string(data)) {
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		if k, v, ok := cut(line, '='); ok && k == key {
			return v
		}
	}
	return ""
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i, c := range s {
		if c == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	lines = append(lines, s[start:])
	return lines
}

func cut(s string, sep byte) (string, string, bool) {
	for i := 0; i < len(s); i++ {
		if s[i] == sep {
			return s[:i], s[i+1:], true
		}
	}
	return s, "", false
}
