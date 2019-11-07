package netauth

import (
	"errors"
)

var (
	// ErrUnknownCache is returned when a cache implementation is
	// requested that has not been registered to the server.
	ErrUnknownCache = errors.New("requested cache is not known")

	// ErrNoCachedToken is returned for a cache miss when asking
	// for a specific owner.  This error should generally prompt
	// the token to be obtained from the server, but may handle
	// terminal failures in the event that a cached token was
	// specifically requested.
	ErrNoCachedToken = errors.New("no cached token for that owner")
)
