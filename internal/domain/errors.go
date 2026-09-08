package domain

import "errors"

var (
	// ErrInvalidUserID indicates an invalid user ID.
	ErrInvalidUserID = errors.New("domain: invalid user id")

	// ErrInvalidSegment indicates an invalid segment.
	ErrInvalidSegment = errors.New("domain: invalid segment")
)
