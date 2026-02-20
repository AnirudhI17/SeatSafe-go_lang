package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"ticketing/backend/internal/config"
	"ticketing/backend/internal/handler"
	"ticketing/backend/internal/repository/postgres"
	"ticketing/backend/internal/router"
	"ticketing/backend/internal/service"
)

func main() {
	// ── Logger ────────────────────────────────────────────────────────────────
	logger, _ := zap.NewProduction()
	defer logger.Sync() //nolint:errcheck

	// ── Config ────────────────────────────────────────────────────────────────
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}
	logger.Info("config loaded", zap.String("env", cfg.App.Env))

	// ── Database ──────────────────────────────────────────────────────────────
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := postgres.NewPool(ctx, cfg.Database)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer pool.Close()
	logger.Info("database connection pool established")

	// ── Repositories ──────────────────────────────────────────────────────────
	userRepo  := postgres.NewUserRepository(pool)
	eventRepo := postgres.NewEventRepository(pool)
	regRepo   := postgres.NewRegistrationRepository(pool)
	ticketRepo := postgres.NewTicketRepository(pool)

	// ── Services ──────────────────────────────────────────────────────────────
	userSvc  := service.NewUserService(userRepo, cfg.JWT)
	eventSvc := service.NewEventService(eventRepo)
	regSvc   := service.NewRegistrationService(regRepo, ticketRepo)

	// ── Handlers ──────────────────────────────────────────────────────────────
	authH  := handler.NewAuthHandler(userSvc)
	eventH := handler.NewEventHandler(eventSvc)
	regH   := handler.NewRegistrationHandler(regSvc)

	// ── Router ────────────────────────────────────────────────────────────────
	engine := router.Setup(cfg, logger, authH, eventH, regH)

	// ── HTTP Server ───────────────────────────────────────────────────────────
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      engine,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine for graceful shutdown
	go func() {
		logger.Info("server starting", zap.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server error", zap.Error(err))
		}
	}()

	// ── Graceful Shutdown ─────────────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")
	shutCtx, shutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutCancel()

	if err := srv.Shutdown(shutCtx); err != nil {
		logger.Fatal("server forced shutdown", zap.Error(err))
	}
	logger.Info("server stopped gracefully")
}
