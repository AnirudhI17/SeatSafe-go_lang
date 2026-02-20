//go:build integration

package integration

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"ticketing/backend/internal/config"
	"ticketing/backend/internal/handler"
	"ticketing/backend/internal/repository/postgres"
	"ticketing/backend/internal/router"
	"ticketing/backend/internal/service"
)

// newTestServer wires up the full HTTP stack (handlers, middleware, router)
// against a real database, suitable for end-to-end auth tests.
func newTestServer(t *testing.T) *gin.Engine {
	t.Helper()

	// Use production-like Gin behaviour but keep logging simple for tests.
	gin.SetMode(gin.TestMode)

	// Minimal config: load from env but override DB + JWT if necessary.
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load: %v", err)
	}

	// Ensure a short-lived but valid JWT expiry for tests.
	if cfg.JWT.ExpiryMinutes == 0 {
		cfg.JWT.ExpiryMinutes = 30
	}
	if cfg.JWT.Secret == "" {
		cfg.JWT.Secret = "test-secret-must-be-long-enough-1234567890"
	}

	// Override DATABASE_URL if provided specifically for tests.
	if testDB := os.Getenv("DATABASE_URL"); testDB != "" {
		cfg.Database.URL = testDB
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)

	pool, err := postgres.NewPool(ctx, cfg.Database)
	if err != nil {
		t.Fatalf("NewPool: %v", err)
	}
	t.Cleanup(pool.Close)

	userRepo := postgres.NewUserRepository(pool)
	eventRepo := postgres.NewEventRepository(pool)
	regRepo := postgres.NewRegistrationRepository(pool)
	ticketRepo := postgres.NewTicketRepository(pool)

	userSvc := service.NewUserService(userRepo, cfg.JWT)
	eventSvc := service.NewEventService(eventRepo)
	regSvc := service.NewRegistrationService(regRepo, ticketRepo)

	authH := handler.NewAuthHandler(userSvc)
	eventH := handler.NewEventHandler(eventSvc)
	regH := handler.NewRegistrationHandler(regSvc)

	logger, _ := zap.NewDevelopment()
	t.Cleanup(func() { _ = logger.Sync() })

	return router.Setup(cfg, logger, authH, eventH, regH)
}

// TestAuth_RegisterLoginAndAccess verifies that a user can register, receive a JWT,
// and use it to access a protected endpoint.
func TestAuth_RegisterLoginAndAccess(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	engine := newTestServer(t)

	// 1. Register a new user
	registerBody := `{"email":"auth_test@example.com","password":"Str0ngPass!","full_name":"Auth Test"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", http.NoBody)
	req.Body = io.NopCloser(strings.NewReader(registerBody))
	req.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created from /auth/register, got %d: %s", w.Code, w.Body.String())
	}

	// 2. Login with same credentials
	loginBody := `{"email":"auth_test@example.com","password":"Str0ngPass!"}`
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", http.NoBody)
	req.Body = io.NopCloser(strings.NewReader(loginBody))
	req.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK from /auth/login, got %d: %s", w.Code, w.Body.String())
	}

	// Extract token from response JSON (simple parse to avoid extra deps).
	var resp struct {
		Success bool `json:"success"`
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("json.Unmarshal login response: %v", err)
	}
	if resp.Data.Token == "" {
		t.Fatal("expected non-empty JWT token in login response")
	}

	// 3. Call protected profile endpoint with Bearer token
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+resp.Data.Token)
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK from /auth/me with valid token, got %d: %s", w.Code, w.Body.String())
	}
}

// TestAuth_ProtectedRoutesRequireToken ensures that protected routes are denied
// when no Authorization header is present.
func TestAuth_ProtectedRoutesRequireToken(t *testing.T) {
	engine := newTestServer(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized without token, got %d", w.Code)
	}
}

// TestAuth_OrganizerRegistrationsEndpointRoleProtected ensures the organiser-only
// registrations endpoint is not accessible without appropriate credentials.
func TestAuth_OrganizerRegistrationsEndpointRoleProtected(t *testing.T) {
	engine := newTestServer(t)

	// Call the organiser endpoint without any token.
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/events/"+uuid.NewString()+"/registrations", nil)
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized for organiser endpoint without token, got %d", w.Code)
	}
}

