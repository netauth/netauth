package netauth

import (
	"errors"
)

var (
	// ErrUnknownCache is returned when a cache implementation is
	// requested that has not been registered to the server.
	ErrUnknownCache = errors.New("The requested cache is not known")
)
