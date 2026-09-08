package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/moritiza/go-challenge/internal/domain"
)

const segmentKeyPrefix = "segment:"

// Storage is the Redis storage client.
type Storage struct {
	client *redis.Client
}

// NewStorage creates a new Redis client.
func NewStorage(addr string, password string, db int) *Storage {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &Storage{client: client}
}

// Ping pings the Redis server.
func (s *Storage) Ping(ctx context.Context) error {
	if err := s.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("ping redis: %w", err)
	}

	return nil
}

// Close closes the Redis client.
func (s *Storage) Close() error {
	return s.client.Close()
}

// StoreMembership stores or refreshes a membership.
func (s *Storage) StoreMembership(ctx context.Context, m *domain.Membership) error {
	err := s.client.ZAdd(ctx, segmentKey(m.Segment), redis.Z{
		Score:  float64(m.ExpiresAt.Unix()),
		Member: m.UserID,
	}).Err()
	if err != nil {
		return fmt.Errorf("store membership: %w", err)
	}

	return nil
}

// CountActiveUsers counts the number of active users in a segment.
func (s *Storage) CountActiveUsers(ctx context.Context, segment string) (int64, error) {
	count, err := s.client.ZCount(
		ctx,
		segmentKey(segment),
		fmt.Sprintf("%d", time.Now().Unix()),
		"+inf",
	).Result()
	if err != nil {
		return 0, fmt.Errorf("count active users: %w", err)
	}

	return count, nil
}

// CleanupExpired cleans up expired memberships from a segment.
func (s *Storage) CleanupExpired(ctx context.Context, segment string) error {
	now := time.Now().Unix()

	if err := s.client.ZRemRangeByScore(
		ctx,
		segmentKey(segment),
		"-inf",
		fmt.Sprintf("%d", now),
	).Err(); err != nil {
		return fmt.Errorf("cleanup expired: %w", err)
	}

	return nil
}

// segmentKey returns the key for a segment.
func segmentKey(segment string) string {
	return segmentKeyPrefix + segment
}
