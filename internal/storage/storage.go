package storage

import (
	"context"

	"github.com/moritiza/go-challenge/internal/domain"
)

// Storage persists memberships and provides segment queries.
type Storage interface {
	// StoreMembership saves or refreshes a membership.
	StoreMembership(ctx context.Context, m *domain.Membership) error

	// CountActiveUsers returns the number of active users in a segment.
	CountActiveUsers(ctx context.Context, segment string) (int64, error)
}
