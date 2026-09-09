package http

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/moritiza/go-challenge/internal/domain"
)

// estimationService is the interface for the estimation service.
type estimationService interface {
	StoreMembership(ctx context.Context, userID, segment string) error
	Estimate(ctx context.Context, segment string) (int64, error)
	Healthy(ctx context.Context) error
}

// Handler is the HTTP handler for the estimation service.
type Handler struct {
	svc    estimationService
	logger *slog.Logger
}

// NewHandler creates a new HTTP handler for the estimation service.
func NewHandler(svc estimationService, logger *slog.Logger) *Handler {
	return &Handler{
		svc:    svc,
		logger: logger,
	}
}

// StoreMembership handles the request to store a membership.
func (h *Handler) StoreMembership(w http.ResponseWriter, r *http.Request) {
	var req storeMembershipRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")

		return
	}

	if err := h.svc.StoreMembership(r.Context(), req.UserID, req.Segment); err != nil {
		h.handleServiceError(w, err, "store membership")

		return
	}

	respondJSON(w, http.StatusCreated, storeMembershipResponse{
		UserID:  req.UserID,
		Segment: req.Segment,
	})
}

// Estimate handles the request to estimate the number of active users in a segment.
func (h *Handler) Estimate(w http.ResponseWriter, r *http.Request) {
	segment := r.PathValue("segment")

	count, err := h.svc.Estimate(r.Context(), segment)
	if err != nil {
		h.handleServiceError(w, err, "estimate segment")

		return
	}

	respondJSON(w, http.StatusOK, estimateResponse{
		Segment: segment,
		Count:   count,
	})
}

// Healthz handles the request for health-check,
// It reports healthy only if the service's dependencies are reachable.
func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.Healthy(r.Context()); err != nil {
		h.logger.Error("health check failed", "error", err)

		respondJSON(w, http.StatusServiceUnavailable, healthzResponse{
			Status: "unhealthy",
		})

		return
	}

	respondJSON(w, http.StatusOK, healthzResponse{
		Status: "ok",
	})
}

// handleServiceError handles service errors.
func (h *Handler) handleServiceError(w http.ResponseWriter, err error, action string) {
	if errors.Is(err, domain.ErrInvalidUserID) ||
		errors.Is(err, domain.ErrInvalidSegment) {
		respondError(w, http.StatusBadRequest, err.Error())

		return
	}

	h.logger.Error(action+" failed", "error", err)

	respondError(w, http.StatusInternalServerError, "internal error")
}

// respondJSON writes a JSON response.
func respondJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(v)
}

// respondError writes an error response.
func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, errorResponse{
		Error: msg,
	})
}
