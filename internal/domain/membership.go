package domain

import "time"

// Membership represents a user's membership in a segment.
// A membership is active for a limited amount of time.
type Membership struct {
	UserID    string
	Segment   string
	CreatedAt time.Time
	ExpiresAt time.Time
}

// NewMembership creates a new membership and sets its expiration time
// based on the given TTL.
func NewMembership(userID string, segment string, ttl time.Duration) (*Membership, error) {
	if userID == "" {
		return nil, ErrInvalidUserID
	}

	if segment == "" {
		return nil, ErrInvalidSegment
	}

	now := time.Now().UTC()

	return &Membership{
		UserID:    userID,
		Segment:   segment,
		CreatedAt: now,
		ExpiresAt: now.Add(ttl),
	}, nil
}

// IsActive returns true if the membership has not expired yet.
func (m *Membership) IsActive(at time.Time) bool {
	return at.Before(m.ExpiresAt)
}
