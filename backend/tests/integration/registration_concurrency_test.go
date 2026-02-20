//go:build integration

// Package integration contains integration tests that run against a real database.
// Build tag `integration` allows normal `go test ./...` to skip these tests.
// Run with: go test -tags=integration -v -race ./tests/integration/...
package integration

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"ticketing/backend/internal/repository"
	"ticketing/backend/internal/repository/postgres"
)

// ─────────────────────────────────────────────────────────────────────────────
// TestConcurrentBooking is the primary concurrency validation test.
//
// Scenario:
//   - Create an event with capacity = 10
//   - Spawn 100 goroutines, each attempting to book 1 seat simultaneously
//   - Assert exactly 10 succeed and 90 fail with ErrEventFull
//   - Assert registered_count in DB is exactly 10 (no overbooking)
//   - Run with -race flag to detect data races in application code
//
// Expected behaviour:
//   SELECT FOR UPDATE serialises all goroutines on the event row.
//   PostgreSQL queues them and grants the lock one at a time.
//   Goroutines 1–10 see count < 10 → succeed.
//   Goroutines 11–100 see count == 10 → return ErrEventFull.
// ─────────────────────────────────────────────────────────────────────────────
func TestConcurrentBooking(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	db := setupTestDB(t, ctx)
	defer cleanupTestDB(t, ctx, db)

	regRepo := postgres.NewRegistrationRepository(db)

	const (
		eventCapacity  = 10
		totalGoroutines = 100
	)

	// ── 1. Seed: organiser + event ───────────────────────────────────────────
	organizerID := seedUser(t, ctx, db, "organizer@test.com")
	eventID     := seedEvent(t, ctx, db, organizerID, eventCapacity)

	// ── 2. Seed: 100 unique users ────────────────────────────────────────────
	userIDs := make([]uuid.UUID, totalGoroutines)
	for i := 0; i < totalGoroutines; i++ {
		userIDs[i] = seedUser(t, ctx, db, fmt.Sprintf("user%d@test.com", i))
	}

	// ── 3. Launch 100 concurrent goroutines at the same instant ─────────────
	var (
		wg         sync.WaitGroup
		successCnt atomic.Int64
		fullCnt    atomic.Int64
		otherErrCnt atomic.Int64
	)

	// Barrier: all goroutines wait until all are created before proceeding
	ready := make(chan struct{})

	for i := 0; i < totalGoroutines; i++ {
		wg.Add(1)
		go func(userID uuid.UUID) {
			defer wg.Done()
			<-ready // Block until all goroutines are ready — maximises concurrency

			bookCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			_, err := regRepo.BookSeat(bookCtx, eventID, userID, 1)
			switch {
			case err == nil:
				successCnt.Add(1)
			case err == repository.ErrEventFull:
				fullCnt.Add(1)
			default:
				otherErrCnt.Add(1)
				t.Logf("unexpected error for user %s: %v", userID, err)
			}
		}(userIDs[i])
	}

	// Fire the starting pistol — unblock all goroutines simultaneously
	close(ready)
	wg.Wait()

	// ── 4. Assertions ─────────────────────────────────────────────────────────
	t.Logf("Results → success: %d | full: %d | other errors: %d",
		successCnt.Load(), fullCnt.Load(), otherErrCnt.Load())

	// CRITICAL: Exactly 10 bookings must succeed
	if got := successCnt.Load(); got != eventCapacity {
		t.Errorf("OVERBOOKING DETECTED: expected %d successes, got %d", eventCapacity, got)
	}

	// CRITICAL: The remaining 90 must report ErrEventFull, not unknown errors
	expectedFull := int64(totalGoroutines - eventCapacity)
	if got := fullCnt.Load(); got != expectedFull {
		t.Errorf("expected %d ErrEventFull rejections, got %d", expectedFull, got)
	}

	// Zero unexpected errors
	if got := otherErrCnt.Load(); got != 0 {
		t.Errorf("expected 0 unexpected errors, got %d", got)
	}

	// CRITICAL: Verify DB state — registered_count must equal exactly 10
	assertRegisteredCount(t, ctx, db, eventID, eventCapacity)

	// CRITICAL: Verify registration table row count
	assertRegistrationRows(t, ctx, db, eventID, eventCapacity)
}

// ─────────────────────────────────────────────────────────────────────────────
// TestConcurrentBooking_SameUser validates the duplicate registration guard
// under concurrency: the same user cannot claim multiple seats even if they
// spam requests simultaneously.
// ─────────────────────────────────────────────────────────────────────────────
func TestConcurrentBooking_SameUser(t *testing.T) {
	ctx := context.Background()
	db  := setupTestDB(t, ctx)
	defer cleanupTestDB(t, ctx, db)

	regRepo := postgres.NewRegistrationRepository(db)

	organizerID := seedUser(t, ctx, db, "organizer2@test.com")
	eventID     := seedEvent(t, ctx, db, organizerID, 50) // plenty of seats
	userID      := seedUser(t, ctx, db, "singleuser@test.com")

	const attempts = 20
	var (
		wg         sync.WaitGroup
		successCnt atomic.Int64
		dupCnt     atomic.Int64
	)
	ready := make(chan struct{})

	for i := 0; i < attempts; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ready
			_, err := regRepo.BookSeat(context.Background(), eventID, userID, 1)
			if err == nil {
				successCnt.Add(1)
			} else if err == repository.ErrAlreadyRegistered {
				dupCnt.Add(1)
			}
		}()
	}
	close(ready)
	wg.Wait()

	// Only 1 of 20 concurrent attempts by the same user must succeed
	if successCnt.Load() != 1 {
		t.Errorf("expected exactly 1 successful booking for the same user, got %d", successCnt.Load())
	}
	if dupCnt.Load() != attempts-1 {
		t.Errorf("expected %d duplicate errors, got %d", attempts-1, dupCnt.Load())
	}

	assertRegisteredCount(t, ctx, db, eventID, 1)
}

// ─────────────────────────────────────────────────────────────────────────────
// TestConcurrentBooking_ExactCapacity tests the edge case where goroutines ==
// capacity — every goroutine should succeed.
// ─────────────────────────────────────────────────────────────────────────────
func TestConcurrentBooking_ExactCapacity(t *testing.T) {
	ctx := context.Background()
	db  := setupTestDB(t, ctx)
	defer cleanupTestDB(t, ctx, db)

	regRepo := postgres.NewRegistrationRepository(db)

	const capacity = 5
	organizerID := seedUser(t, ctx, db, "organizer3@test.com")
	eventID     := seedEvent(t, ctx, db, organizerID, capacity)
	userIDs     := make([]uuid.UUID, capacity)
	for i := 0; i < capacity; i++ {
		userIDs[i] = seedUser(t, ctx, db, fmt.Sprintf("exact%d@test.com", i))
	}

	var (
		wg         sync.WaitGroup
		successCnt atomic.Int64
	)
	ready := make(chan struct{})
	for _, uid := range userIDs {
		wg.Add(1)
		go func(u uuid.UUID) {
			defer wg.Done()
			<-ready
			_, err := regRepo.BookSeat(context.Background(), eventID, u, 1)
			if err == nil {
				successCnt.Add(1)
			}
		}(uid)
	}
	close(ready)
	wg.Wait()

	// All 5 must succeed — no false rejections at exact capacity
	if successCnt.Load() != capacity {
		t.Errorf("expected all %d bookings to succeed, got %d", capacity, successCnt.Load())
	}
	assertRegisteredCount(t, ctx, db, eventID, capacity)
}
