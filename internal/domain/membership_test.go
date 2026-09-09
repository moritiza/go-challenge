package domain

import (
	"errors"
	"testing"
	"time"
)

func TestNewMembership(t *testing.T) {
	ttl := 14 * 24 * time.Hour

	membership, err := NewMembership(
		"user-1",
		"sports",
		ttl,
	)
	if err != nil {
		t.Fatalf("NewMembership() error = %v", err)
	}

	if membership.UserID != "user-1" {
		t.Fatalf("UserID = %q, want %q",
			membership.UserID,
			"user-1",
		)
	}

	if membership.Segment != "sports" {
		t.Fatalf("Segment = %q, want %q",
			membership.Segment,
			"sports",
		)
	}

	if membership.CreatedAt.IsZero() {
		t.Fatal("CreatedAt should not be zero")
	}

	expectedExpiry := membership.CreatedAt.Add(ttl)

	if !membership.ExpiresAt.Equal(expectedExpiry) {
		t.Fatalf(
			"ExpiresAt = %v, want %v",
			membership.ExpiresAt,
			expectedExpiry,
		)
	}
}

func TestNewMembership_InvalidUserID(t *testing.T) {
	_, err := NewMembership(
		"",
		"sports",
		time.Hour,
	)

	if !errors.Is(err, ErrInvalidUserID) {
		t.Fatalf(
			"NewMembership() error = %v, want %v",
			err,
			ErrInvalidUserID,
		)
	}
}

func TestNewMembership_InvalidSegment(t *testing.T) {
	_, err := NewMembership(
		"user-1",
		"",
		time.Hour,
	)

	if !errors.Is(err, ErrInvalidSegment) {
		t.Fatalf(
			"NewMembership() error = %v, want %v",
			err,
			ErrInvalidSegment,
		)
	}
}

func TestMembership_IsActive(t *testing.T) {
	now := time.Now().UTC()

	membership := Membership{
		UserID:    "user-1",
		Segment:   "sports",
		CreatedAt: now,
		ExpiresAt: now.Add(time.Hour),
	}

	tests := []struct {
		name string
		at   time.Time
		want bool
	}{
		{
			name: "before expiration",
			at:   now.Add(30 * time.Minute),
			want: true,
		},
		{
			name: "at expiration",
			at:   membership.ExpiresAt,
			want: false,
		},
		{
			name: "after expiration",
			at:   now.Add(2 * time.Hour),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := membership.IsActive(tt.at); got != tt.want {
				t.Fatalf(
					"IsActive() = %v, want %v",
					got,
					tt.want,
				)
			}
		})
	}
}
