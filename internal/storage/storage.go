package storage

import (
	"context"

	"github.com/moritiza/go-challenge/internal/domain"
)

// Storage persists memberships and provides segment queries.
type Storage interface {
	// Ping checks that the storage backend is reachable.
	Ping(ctx context.Context) error

	// StoreMembership saves or refreshes a membership.
	StoreMembership(ctx context.Context, m *domain.Membership) error

	// CountActiveUsers returns the number of active users in a segment.
	CountActiveUsers(ctx context.Context, segment string) (int64, error)

	// CleanupExpired removes expired memberships from a segment.
	// It is intended to run periodically as a background cleanup task.
	CleanupExpired(ctx context.Context, segment string) error
}
