// Package memory provides an in-memory implementation of the storage interface.
// It is used for testing and development.
package memory

import (
	"context"
	"sync"
	"time"

	"github.com/moritiza/go-challenge/internal/domain"
)

// Storage is the in-memory storage client.
type Storage struct {
	mu sync.RWMutex

	// segments maps segment -> userID -> expiration time.
	segments map[string]map[string]time.Time
}

// NewStorage creates a new in-memory storage.
func NewStorage() *Storage {
	return &Storage{
		segments: make(map[string]map[string]time.Time),
	}
}

// Ping always succeeds for the in-memory storage,
// there's no external dependency to check.
func (s *Storage) Ping(_ context.Context) error {
	return nil
}

// StoreMembership stores or refreshes a membership.
func (s *Storage) StoreMembership(_ context.Context, m *domain.Membership) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.segments[m.Segment] == nil {
		s.segments[m.Segment] = make(map[string]time.Time)
	}

	s.segments[m.Segment][m.UserID] = m.ExpiresAt

	return nil
}

// CountActiveUsers counts the number of active users in a segment.
func (s *Storage) CountActiveUsers(_ context.Context, segment string) (int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	now := time.Now()
	var count int64

	for _, expiresAt := range s.segments[segment] {
		if now.Before(expiresAt) {
			count++
		}
	}

	return count, nil
}

// CleanupExpired cleans up expired memberships from a segment.
func (s *Storage) CleanupExpired(_ context.Context, segment string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()

	for userID, expiresAt := range s.segments[segment] {
		if !now.Before(expiresAt) {
			delete(s.segments[segment], userID)
		}
	}

	return nil
}
