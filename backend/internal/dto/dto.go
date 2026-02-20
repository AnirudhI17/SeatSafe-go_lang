package dto

import "time"

// ── Event DTOs ────────────────────────────────────────────────────────────────

type CreateEventRequest struct {
	Title       string    `json:"title"        binding:"required,min=3,max=255"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	IsOnline    bool      `json:"is_online"`
	OnlineURL   string    `json:"online_url"`
	StartsAt    time.Time `json:"starts_at"`
	EndsAt      time.Time `json:"ends_at"`
	Capacity    int       `json:"capacity"     binding:"required,min=1"`
	PriceCents  int       `json:"price_cents"  binding:"min=0"`
	BannerURL   string    `json:"banner_url"`
	
	// Alternative field names for backward compatibility
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Price       float64   `json:"price"`
}

type UpdateEventRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Capacity    int       `json:"capacity"    binding:"omitempty,min=1"`
	StartsAt    time.Time `json:"starts_at"`
	EndsAt      time.Time `json:"ends_at"`
}

// ── Registration DTOs ─────────────────────────────────────────────────────────

type BookEventRequest struct {
	Quantity int    `json:"quantity" binding:"min=1,max=10"`
	Notes    string `json:"notes"`
}

// ── User DTOs ─────────────────────────────────────────────────────────────────

type RegisterUserRequest struct {
	Email    string `json:"email"     binding:"required,email"`
	Password string `json:"password"  binding:"required,min=8"`
	FullName string `json:"full_name" binding:"required,min=2"`
	Role     string `json:"role"      binding:"required,oneof=attendee organizer admin"`
}

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  UserProfile `json:"user"`
}

type UserProfile struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

// ── Common ────────────────────────────────────────────────────────────────────

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

func OK(data interface{}) APIResponse {
	return APIResponse{Success: true, Data: data}
}

func Err(msg string) APIResponse {
	return APIResponse{Success: false, Error: msg}
}
