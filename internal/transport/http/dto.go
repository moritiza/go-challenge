package http

// storeMembershipRequest is the request body for storing a membership.
type storeMembershipRequest struct {
	UserID  string `json:"user_id"`
	Segment string `json:"segment"`
}

// storeMembershipResponse is the response body for storing a membership.
type storeMembershipResponse struct {
	UserID  string `json:"user_id"`
	Segment string `json:"segment"`
}

// estimateResponse is the response body for estimating the number of active users in a segment.
type estimateResponse struct {
	Segment string `json:"segment"`
	Count   int64  `json:"count"`
}

// errorResponse is the error response body.
type errorResponse struct {
	Error string `json:"error"`
}

// healthzResponse is health-check response body.
type healthzResponse struct {
	Status string `json:"status"`
}
