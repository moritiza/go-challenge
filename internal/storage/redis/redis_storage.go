package redis

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/moritiza/go-challenge/internal/domain"
)

const (
	segmentKeyPrefix = "segment:"

	// scanCount is the COUNT hint for Redis SCAN.
	// It's a hint, not a hard limit.
	scanCount = 100
)

// Storage is the Redis storage client.
type Storage struct {
	client    *redis.Client
	batchSize int64
}

// NewStorage creates a new Redis client.
func NewStorage(addr string, password string, db int, batchSize int64) *Storage {
	if batchSize <= 0 {
		batchSize = 1000
	}

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &Storage{
		client:    client,
		batchSize: batchSize,
	}
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
	key := segmentKey(segment)
	now := fmt.Sprintf("%d", time.Now().Unix())

	for {
		expired, err := s.client.ZRangeArgs(ctx, redis.ZRangeArgs{
			Key:     key,
			Start:   "-inf",
			Stop:    now,
			ByScore: true,
			Offset:  0,
			Count:   s.batchSize,
		}).Result()
		if err != nil {
			return fmt.Errorf("find expired members: %w", err)
		}

		if len(expired) == 0 {
			return nil
		}

		if err := s.client.ZRem(
			ctx,
			key,
			toInterfaceSlice(expired)...,
		).Err(); err != nil {
			return fmt.Errorf("remove expired members: %w", err)
		}

		if int64(len(expired)) < s.batchSize {
			return nil
		}
	}
}

// ListSegments returns all known segment names.
func (s *Storage) ListSegments(ctx context.Context) ([]string, error) {
	var segments []string
	var cursor uint64

	for {
		keys, next, err := s.client.Scan(
			ctx,
			cursor,
			segmentKeyPrefix+"*",
			scanCount,
		).Result()
		if err != nil {
			return nil, fmt.Errorf("scan segment keys: %w", err)
		}

		for _, key := range keys {
			segments = append(
				segments,
				strings.TrimPrefix(key, segmentKeyPrefix),
			)
		}

		cursor = next
		if cursor == 0 {
			return segments, nil
		}
	}
}

// segmentKey returns the Redis key for a segment.
func segmentKey(segment string) string {
	return segmentKeyPrefix + segment
}

// toInterfaceSlice converts strings to a slice of interfaces.
func toInterfaceSlice(values []string) []interface{} {
	result := make([]interface{}, len(values))

	for i, value := range values {
		result[i] = value
	}

	return result
}
