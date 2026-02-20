//go:build integration

package integration

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/google/uuid"
	"ticketing/backend/internal/repository/postgres"
)

// BenchmarkConcurrentBooking measures booking throughput under concurrency.
// Run: go test -tags=integration -bench=. -benchtime=10s ./tests/integration/...
func BenchmarkConcurrentBooking(b *testing.B) {
	ctx := context.Background()
	db  := setupTestDB((*testing.T)(nil), ctx)
	defer db.Close()

	regRepo := postgres.NewRegistrationRepository(db)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// Each benchmark iteration gets a fresh event with 50 seats
		organizerID := uuid.New()
		// Seed organizer
		db.Exec(ctx, `INSERT INTO users (id, email, password_hash, full_name, role) VALUES ($1,$2,'$2a$12$x','bench org','organizer') ON CONFLICT DO NOTHING`,
			organizerID, fmt.Sprintf("benchorg%d@test.com", i))
		eventID := uuid.New()
		db.Exec(ctx, `INSERT INTO events (id, organizer_id, title, starts_at, ends_at, capacity, price_cents, currency, status)
			VALUES ($1,$2,'Bench',$3,$4,50,0,'USD','published')`,
			eventID, organizerID, "2027-01-01T10:00:00Z", "2027-01-01T18:00:00Z")

		userIDs := make([]uuid.UUID, 50)
		for j := 0; j < 50; j++ {
			uid := uuid.New()
			db.Exec(ctx, `INSERT INTO users (id, email, password_hash, full_name, role) VALUES ($1,$2,'$2a$12$x','bench user','attendee')`,
				uid, fmt.Sprintf("bench%d-%d@test.com", i, j))
			userIDs[j] = uid
		}
		b.StartTimer()

		var wg sync.WaitGroup
		var successCnt atomic.Int64
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

		b.ReportMetric(float64(successCnt.Load()), "bookings/iter")
	}
}
