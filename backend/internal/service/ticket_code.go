package service

import (
	"crypto/rand"
	"fmt"
)

// generateTicketCode creates a human-readable unique ticket code.
// Format: TKT-XXXX-XXXX (e.g. TKT-A3F9-2KXP)
func generateTicketCode() string {
	b := make([]byte, 4)
	rand.Read(b) //nolint:errcheck
	part1 := fmt.Sprintf("%X", b[:2])
	part2 := fmt.Sprintf("%X", b[2:])
	return fmt.Sprintf("TKT-%s-%s", part1, part2)
}
