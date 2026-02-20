package repository

import (
	"errors"
	"fmt"
)

// Sentinel errors — used by repositories and handled by the service layer.
var (
	// ErrNotFound is returned when a record does not exist.
	ErrNotFound = errors.New("record not found")

	// ErrEventFull is returned when an event has no remaining capacity.
	// This is the key error surfaced during concurrent booking races.
	ErrEventFull = errors.New("event is fully booked")

	// ErrAlreadyRegistered is returned when a user tries to register for an event twice.
	ErrAlreadyRegistered = errors.New("already registered for this event")

	// ErrEventNotPublished is returned when booking an event that isn't published.
	ErrEventNotPublished = errors.New("event is not available for registration")

	// ErrEventCancelled is returned when the event has been cancelled.
	ErrEventCancelled = errors.New("event has been cancelled")

	// ErrUnauthorised is returned when a user lacks permission for an operation.
	ErrUnauthorised = errors.New("not authorised to perform this action")

	// ErrDuplicateEmail is returned on unique email violation.
	ErrDuplicateEmail = errors.New("email address is already registered")
)

// ConflictError wraps a conflict with additional context.
type ConflictError struct {
	Message string
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("conflict: %s", e.Message)
}
