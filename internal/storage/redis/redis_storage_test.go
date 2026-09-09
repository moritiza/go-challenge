package redis

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/moritiza/go-challenge/internal/domain"
)

func TestStoreMembership(t *testing.T) {
	store, client := newTestStorage(t)
	segment := testSegment(t)

	membership, err := domain.NewMembership(
		"user-1",
		segment,
		time.Hour,
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := store.StoreMembership(context.Background(), membership); err != nil {
		t.Fatalf("StoreMembership() error = %v", err)
	}

	count, err := client.ZCard(context.Background(), segmentKey(segment)).Result()
	if err != nil {
		t.Fatalf("ZCard() error = %v", err)
	}

	if count != 1 {
		t.Fatalf("ZCard() = %d, want 1", count)
	}
}

func TestStoreMembership_RefreshesExpiration(t *testing.T) {
	store, client := newTestStorage(t)
	segment := testSegment(t)

	first, err := domain.NewMembership("user-1", segment, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	second, err := domain.NewMembership("user-1", segment, 2*time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	if err := store.StoreMembership(ctx, first); err != nil {
		t.Fatalf("first StoreMembership() error = %v", err)
	}

	if err := store.StoreMembership(ctx, second); err != nil {
		t.Fatalf("second StoreMembership() error = %v", err)
	}

	count, err := client.ZCard(ctx, segmentKey(segment)).Result()
	if err != nil {
		t.Fatalf("ZCard() error = %v", err)
	}

	if count != 1 {
		t.Fatalf("ZCard() = %d, want 1", count)
	}

	score, err := client.ZScore(ctx, segmentKey(segment), "user-1").Result()
	if err != nil {
		t.Fatalf("ZScore() error = %v", err)
	}

	if int64(score) != second.ExpiresAt.Unix() {
		t.Fatalf("ZScore() = %d, want %d",
			int64(score),
			second.ExpiresAt.Unix(),
		)
	}
}

func TestCountActiveUsers(t *testing.T) {
	store, _ := newTestStorage(t)
	segment := testSegment(t)
	ctx := context.Background()

	memberships := []struct {
		userID string
		ttl    time.Duration
	}{
		{"user-1", time.Hour},
		{"user-2", 2 * time.Hour},
		{"user-3", -time.Hour},
	}

	for _, item := range memberships {
		membership, err := domain.NewMembership(
			item.userID,
			segment,
			item.ttl,
		)
		if err != nil {
			t.Fatal(err)
		}

		if err := store.StoreMembership(ctx, membership); err != nil {
			t.Fatalf("StoreMembership() error = %v", err)
		}
	}

	count, err := store.CountActiveUsers(ctx, segment)
	if err != nil {
		t.Fatalf("CountActiveUsers() error = %v", err)
	}

	if count != 2 {
		t.Fatalf("CountActiveUsers() = %d, want 2", count)
	}
}

func TestCleanupExpired(t *testing.T) {
	store, client := newTestStorage(t)
	segment := testSegment(t)
	ctx := context.Background()

	memberships := []struct {
		userID string
		ttl    time.Duration
	}{
		{"user-1", -time.Hour},
		{"user-2", -2 * time.Hour},
		{"user-3", time.Hour},
	}

	for _, item := range memberships {
		membership, err := domain.NewMembership(
			item.userID,
			segment,
			item.ttl,
		)
		if err != nil {
			t.Fatal(err)
		}

		if err := store.StoreMembership(ctx, membership); err != nil {
			t.Fatalf("StoreMembership() error = %v", err)
		}
	}

	if err := store.CleanupExpired(ctx, segment); err != nil {
		t.Fatalf("CleanupExpired() error = %v", err)
	}

	count, err := client.ZCard(ctx, segmentKey(segment)).Result()
	if err != nil {
		t.Fatalf("ZCard() error = %v", err)
	}

	if count != 1 {
		t.Fatalf("ZCard() = %d, want 1", count)
	}

	if _, err := client.ZScore(
		ctx,
		segmentKey(segment),
		"user-3",
	).Result(); err != nil {
		t.Fatalf("active membership was removed: %v", err)
	}
}

func TestListSegments(t *testing.T) {
	store, _ := newTestStorage(t)
	ctx := context.Background()

	segments := []string{
		testSegment(t) + "-a",
		testSegment(t) + "-b",
		testSegment(t) + "-c",
	}

	for _, segment := range segments {
		membership, err := domain.NewMembership(
			"user-1",
			segment,
			time.Hour,
		)
		if err != nil {
			t.Fatal(err)
		}

		if err := store.StoreMembership(ctx, membership); err != nil {
			t.Fatalf("StoreMembership() error = %v", err)
		}
	}

	got, err := store.ListSegments(ctx)
	if err != nil {
		t.Fatalf("ListSegments() error = %v", err)
	}

	gotSet := make(map[string]bool, len(got))
	for _, segment := range got {
		gotSet[segment] = true
	}

	for _, segment := range segments {
		if !gotSet[segment] {
			t.Errorf("ListSegments() missing segment %q", segment)
		}
	}
}

func TestPing(t *testing.T) {
	store, _ := newTestStorage(t)

	if err := store.Ping(context.Background()); err != nil {
		t.Fatalf("Ping() error = %v", err)
	}
}

func newTestStorage(t *testing.T) (*Storage, *redis.Client) {
	t.Helper()

	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		t.Skipf("Redis is not available at %s", addr)
	}

	t.Cleanup(func() {
		_ = client.Close()
	})

	return &Storage{
		client:    client,
		batchSize: 2,
	}, client
}

func testSegment(t *testing.T) string {
	t.Helper()

	return "test:" + strings.ReplaceAll(t.Name(), "/", ":")
}
