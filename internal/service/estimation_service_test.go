package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/moritiza/go-challenge/internal/domain"
	"github.com/moritiza/go-challenge/internal/storage/memory"
)

func TestStoreMembership(t *testing.T) {
	store := memory.NewStorage()
	svc := NewEstimationService(store, time.Hour)

	err := svc.StoreMembership(
		context.Background(),
		"user-1",
		"sports",
	)
	if err != nil {
		t.Fatalf("StoreMembership() error = %v", err)
	}

	count, err := svc.Estimate(context.Background(), "sports")
	if err != nil {
		t.Fatalf("Estimate() error = %v", err)
	}

	if count != 1 {
		t.Fatalf("Estimate() = %d, want 1", count)
	}
}

func TestStoreMembership_RefreshesMembership(t *testing.T) {
	store := memory.NewStorage()
	svc := NewEstimationService(store, time.Hour)

	ctx := context.Background()

	if err := svc.StoreMembership(ctx, "user-1", "sports"); err != nil {
		t.Fatalf("first StoreMembership() error = %v", err)
	}

	if err := svc.StoreMembership(ctx, "user-1", "sports"); err != nil {
		t.Fatalf("second StoreMembership() error = %v", err)
	}

	count, err := svc.Estimate(ctx, "sports")
	if err != nil {
		t.Fatalf("Estimate() error = %v", err)
	}

	if count != 1 {
		t.Fatalf("Estimate() = %d, want 1", count)
	}
}

func TestStoreMembership_InvalidUserID(t *testing.T) {
	store := memory.NewStorage()
	svc := NewEstimationService(store, time.Hour)

	err := svc.StoreMembership(
		context.Background(),
		"",
		"sports",
	)

	if !errors.Is(err, domain.ErrInvalidUserID) {
		t.Fatalf("StoreMembership() error = %v, want %v",
			err,
			domain.ErrInvalidUserID,
		)
	}
}

func TestStoreMembership_InvalidSegment(t *testing.T) {
	store := memory.NewStorage()
	svc := NewEstimationService(store, time.Hour)

	err := svc.StoreMembership(
		context.Background(),
		"user-1",
		"",
	)

	if !errors.Is(err, domain.ErrInvalidSegment) {
		t.Fatalf("StoreMembership() error = %v, want %v",
			err,
			domain.ErrInvalidSegment,
		)
	}
}

func TestEstimate(t *testing.T) {
	store := memory.NewStorage()
	svc := NewEstimationService(store, time.Hour)

	ctx := context.Background()

	memberships := []struct {
		userID  string
		segment string
	}{
		{"user-1", "sports"},
		{"user-2", "sports"},
		{"user-3", "sports"},
		{"user-4", "music"},
	}

	for _, membership := range memberships {
		if err := svc.StoreMembership(
			ctx,
			membership.userID,
			membership.segment,
		); err != nil {
			t.Fatalf("StoreMembership() error = %v", err)
		}
	}

	count, err := svc.Estimate(ctx, "sports")
	if err != nil {
		t.Fatalf("Estimate() error = %v", err)
	}

	if count != 3 {
		t.Fatalf("Estimate() = %d, want 3", count)
	}
}

func TestEstimate_EmptySegment(t *testing.T) {
	store := memory.NewStorage()
	svc := NewEstimationService(store, time.Hour)

	_, err := svc.Estimate(context.Background(), "")

	if !errors.Is(err, domain.ErrInvalidSegment) {
		t.Fatalf("Estimate() error = %v, want %v",
			err,
			domain.ErrInvalidSegment,
		)
	}
}

func TestEstimate_UnknownSegment(t *testing.T) {
	store := memory.NewStorage()
	svc := NewEstimationService(store, time.Hour)

	count, err := svc.Estimate(
		context.Background(),
		"unknown",
	)
	if err != nil {
		t.Fatalf("Estimate() error = %v", err)
	}

	if count != 0 {
		t.Fatalf("Estimate() = %d, want 0", count)
	}
}

func TestCleanupExpired(t *testing.T) {
	store := memory.NewStorage()
	svc := NewEstimationService(store, time.Hour)

	ctx := context.Background()

	if err := svc.StoreMembership(ctx, "user-1", "sports"); err != nil {
		t.Fatalf("StoreMembership() error = %v", err)
	}

	if err := svc.CleanupExpired(ctx, "sports"); err != nil {
		t.Fatalf("CleanupExpired() error = %v", err)
	}

	count, err := svc.Estimate(ctx, "sports")
	if err != nil {
		t.Fatalf("Estimate() error = %v", err)
	}

	if count != 1 {
		t.Fatalf("Estimate() = %d, want 1", count)
	}
}

func TestCleanupExpired_EmptySegment(t *testing.T) {
	store := memory.NewStorage()
	svc := NewEstimationService(store, time.Hour)

	err := svc.CleanupExpired(
		context.Background(),
		"",
	)

	if !errors.Is(err, domain.ErrInvalidSegment) {
		t.Fatalf("CleanupExpired() error = %v, want %v",
			err,
			domain.ErrInvalidSegment,
		)
	}
}

func TestCleanupAllExpired(t *testing.T) {
	store := memory.NewStorage()
	svc := NewEstimationService(store, time.Hour)

	ctx := context.Background()

	if err := svc.StoreMembership(ctx, "user-1", "sports"); err != nil {
		t.Fatalf("StoreMembership() error = %v", err)
	}

	if err := svc.StoreMembership(ctx, "user-2", "music"); err != nil {
		t.Fatalf("StoreMembership() error = %v", err)
	}

	if err := svc.CleanupAllExpired(ctx); err != nil {
		t.Fatalf("CleanupAllExpired() error = %v", err)
	}

	sportsCount, err := svc.Estimate(ctx, "sports")
	if err != nil {
		t.Fatalf("Estimate() error = %v", err)
	}

	musicCount, err := svc.Estimate(ctx, "music")
	if err != nil {
		t.Fatalf("Estimate() error = %v", err)
	}

	if sportsCount != 1 {
		t.Fatalf("sports count = %d, want 1", sportsCount)
	}

	if musicCount != 1 {
		t.Fatalf("music count = %d, want 1", musicCount)
	}
}

func TestHealthy(t *testing.T) {
	store := memory.NewStorage()
	svc := NewEstimationService(store, time.Hour)

	if err := svc.Healthy(context.Background()); err != nil {
		t.Fatalf("Healthy() error = %v", err)
	}
}
