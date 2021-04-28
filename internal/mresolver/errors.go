package mresolver

import "errors"

var (
	// ErrInsufficientKnowledge is returned when a resolution is
	// requested that depends on information that the resolver
	// doesn't currently have.
	ErrInsufficientKnowledge = errors.New("insufficient knowledge to satisfy request")
)
