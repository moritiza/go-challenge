package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/moritiza/go-challenge/internal/domain"
	"github.com/moritiza/go-challenge/internal/storage"
)

// EstimationService is the estimation service.
type EstimationService struct {
	store storage.Storage
	ttl   time.Duration
}

// NewEstimationService creates a new estimation service.
func NewEstimationService(store storage.Storage, ttl time.Duration) *EstimationService {
	return &EstimationService{
		store: store,
		ttl:   ttl,
	}
}

// StoreMembership stores or refreshes a user's membership.
func (s *EstimationService) StoreMembership(ctx context.Context, userID string, segment string) error {
	membership, err := domain.NewMembership(userID, segment, s.ttl)
	if err != nil {
		return fmt.Errorf("build membership: %w", err)
	}

	if err := s.store.StoreMembership(ctx, membership); err != nil {
		return fmt.Errorf("store membership: %w", err)
	}

	return nil
}

// Estimate counts the number of active users in a segment.
func (s *EstimationService) Estimate(ctx context.Context, segment string) (int64, error) {
	if segment == "" {
		return 0, domain.ErrInvalidSegment
	}

	count, err := s.store.CountActiveUsers(ctx, segment)
	if err != nil {
		return 0, fmt.Errorf("count active users: %w", err)
	}

	return count, nil
}

// CleanupExpired cleans up expired memberships from a segment.
func (s *EstimationService) CleanupExpired(ctx context.Context, segment string) error {
	if segment == "" {
		return domain.ErrInvalidSegment
	}

	if err := s.store.CleanupExpired(ctx, segment); err != nil {
		return fmt.Errorf("cleanup expired: %w", err)
	}

	return nil
}

// CleanupAllExpired cleans up expired memberships from all segments.
func (s *EstimationService) CleanupAllExpired(ctx context.Context) error {
	segments, err := s.store.ListSegments(ctx)
	if err != nil {
		return fmt.Errorf("list segments: %w", err)
	}

	var errs []error

	for _, segment := range segments {
		if err := s.CleanupExpired(ctx, segment); err != nil {
			errs = append(
				errs,
				fmt.Errorf("cleanup segment %q: %w", segment, err),
			)
		}
	}

	return errors.Join(errs...)
}

// Healthy reports whether the service's dependencies are reachable.
func (s *EstimationService) Healthy(ctx context.Context) error {
	return s.store.Ping(ctx)
}
