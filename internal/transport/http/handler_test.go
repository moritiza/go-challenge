package http

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/moritiza/go-challenge/internal/service"
	"github.com/moritiza/go-challenge/internal/storage/memory"
)

func newTestHandler() *Handler {
	store := memory.NewStorage()
	svc := service.NewEstimationService(store, time.Hour)

	logger := slog.New(slog.NewTextHandler(
		httptest.NewRecorder(),
		nil,
	))

	return NewHandler(svc, logger)
}

func TestStoreMembership(t *testing.T) {
	handler := newTestHandler()

	req := httptest.NewRequest(
		http.MethodPost,
		"/v1/memberships",
		strings.NewReader(`{"user_id":"user-1","segment":"sports"}`),
	)
	rec := httptest.NewRecorder()

	handler.StoreMembership(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d",
			rec.Code,
			http.StatusCreated,
		)
	}

	want := `{"user_id":"user-1","segment":"sports"}`

	if strings.TrimSpace(rec.Body.String()) != want {
		t.Fatalf("body = %q, want %q",
			strings.TrimSpace(rec.Body.String()),
			want,
		)
	}
}

func TestStoreMembership_InvalidJSON(t *testing.T) {
	handler := newTestHandler()

	req := httptest.NewRequest(
		http.MethodPost,
		"/v1/memberships",
		strings.NewReader(`{"user_id":`),
	)
	rec := httptest.NewRecorder()

	handler.StoreMembership(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d",
			rec.Code,
			http.StatusBadRequest,
		)
	}
}

func TestStoreMembership_InvalidInput(t *testing.T) {
	handler := newTestHandler()

	req := httptest.NewRequest(
		http.MethodPost,
		"/v1/memberships",
		strings.NewReader(`{"user_id":"","segment":"sports"}`),
	)
	rec := httptest.NewRecorder()

	handler.StoreMembership(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d",
			rec.Code,
			http.StatusBadRequest,
		)
	}
}

func TestEstimate(t *testing.T) {
	handler := newTestHandler()

	handler.StoreMembership(
		httptest.NewRecorder(),
		httptest.NewRequest(
			http.MethodPost,
			"/v1/memberships",
			strings.NewReader(`{"user_id":"user-1","segment":"sports"}`),
		),
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/v1/segments/sports/estimate",
		nil,
	)
	req.SetPathValue("segment", "sports")

	rec := httptest.NewRecorder()

	handler.Estimate(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d",
			rec.Code,
			http.StatusOK,
		)
	}

	want := `{"segment":"sports","count":1}`

	if strings.TrimSpace(rec.Body.String()) != want {
		t.Fatalf("body = %q, want %q",
			strings.TrimSpace(rec.Body.String()),
			want,
		)
	}
}

func TestEstimate_InvalidSegment(t *testing.T) {
	handler := newTestHandler()

	req := httptest.NewRequest(
		http.MethodGet,
		"/v1/segments//estimate",
		nil,
	)
	req.SetPathValue("segment", "")

	rec := httptest.NewRecorder()

	handler.Estimate(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d",
			rec.Code,
			http.StatusBadRequest,
		)
	}
}

func TestHealthz(t *testing.T) {
	handler := newTestHandler()

	req := httptest.NewRequest(
		http.MethodGet,
		"/healthz",
		nil,
	)
	rec := httptest.NewRecorder()

	handler.Healthz(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d",
			rec.Code,
			http.StatusOK,
		)
	}
}

func TestRespondError(t *testing.T) {
	rec := httptest.NewRecorder()

	respondError(
		rec,
		http.StatusBadRequest,
		"invalid request",
	)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d",
			rec.Code,
			http.StatusBadRequest,
		)
	}
}
