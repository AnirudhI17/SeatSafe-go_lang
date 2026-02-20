package domain

import (
	"time"

	"github.com/google/uuid"
)

// UserRole defines RBAC roles matching the DB enum.
type UserRole string

const (
	RoleAttendee  UserRole = "attendee"
	RoleOrganizer UserRole = "organizer"
	RoleAdmin     UserRole = "admin"
)

// User represents a platform user.
type User struct {
	ID            uuid.UUID `json:"id"`
	Email         string    `json:"email"`
	PasswordHash  string    `json:"-"` // never serialised to JSON
	FullName      string    `json:"full_name"`
	Role          UserRole  `json:"role"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
