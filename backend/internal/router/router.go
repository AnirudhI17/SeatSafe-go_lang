package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"ticketing/backend/internal/config"
	"ticketing/backend/internal/handler"
	"ticketing/backend/internal/middleware"
)

// Setup wires all routes and returns the configured Gin engine.
func Setup(
	cfg *config.Config,
	logger *zap.Logger,
	authH  *handler.AuthHandler,
	eventH *handler.EventHandler,
	regH   *handler.RegistrationHandler,
) *gin.Engine {
	if cfg.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger(logger))
	r.Use(middleware.CORS(cfg.CORS.AllowedOrigins))
	r.Use(middleware.ErrorHandler(logger))

	// ── Health Check ──────────────────────────────────────────────────────────
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "ticketing-api",
			"version": "1.0.0",
		})
	})

	// ── API v1 ────────────────────────────────────────────────────────────────
	v1 := r.Group("/api/v1")

	// ── Auth (public) ─────────────────────────────────────────────────────────
	auth := v1.Group("/auth")
	{
		auth.POST("/register", authH.Register)
		auth.POST("/login",    authH.Login)
	}

	// ── Events (public reads) ─────────────────────────────────────────────────
	v1.GET("/events",     eventH.ListEvents)
	v1.GET("/events/:id", eventH.GetEvent)

	// ── Authenticated routes ──────────────────────────────────────────────────
	secure := v1.Group("/")
	secure.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
	{
		// Auth (protected)
		secure.GET("/auth/me", authH.Profile)

		// Event management — organizer/admin only
		organizer := secure.Group("/")
		organizer.Use(middleware.RequireRole("organizer", "admin"))
		{
			organizer.POST("/events",                     eventH.CreateEvent)
			organizer.PATCH("/events/:id/publish",        eventH.PublishEvent)
			organizer.GET("/events/:id/registrations",    regH.ListEventRegistrations)
		}

		// Registrations — any authenticated user (attendee, organizer, admin)
		secure.POST("/events/:id/register",  regH.BookEvent)
		secure.GET("/registrations/me",      regH.MyRegistrations)
		secure.DELETE("/registrations/:id",  regH.CancelRegistration)
		secure.GET("/tickets/me",            regH.MyTickets)
	}

	return r
}
