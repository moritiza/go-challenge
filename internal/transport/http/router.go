package http

import "net/http"

// NewRouter builds the HTTP routes for the estimation service.
func NewRouter(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", h.Healthz)

	mux.HandleFunc("POST /v1/memberships", h.StoreMembership)
	mux.HandleFunc("GET /v1/segments/{segment}/estimate", h.Estimate)

	return mux
}
