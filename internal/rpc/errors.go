package rpc

import (
	"errors"
)

var (
	// ErrRequestorUnqualified is returned when a caller has
	// attempted to perform some action that requires
	// authorization and the caller is either not authorized, was
	// unable to present a token, or the token did not contain
	// sufficient capabilities.
	ErrRequestorUnqualified = errors.New("the requestor is not qualified to perform that action")

	// ErrMalformedRequest is returned when a caller makes some
	// request to the server but has failed to provide a complete
	// request, or has provided a request that is in conflict with
	// itself.
	ErrMalformedRequest = errors.New("the request is malformed and cannot be processed")

	// ErrReadOnly is used to signal to the client that they have
	// requested a mutating change from a server that is unable to
	// comply due to being in read-only mode.
	ErrReadOnly = errors.New("this server is in read only mode")

	// ErrInternalError is a catchall for errors that are
	// otherwise unidentified and unrecoverable in the server.
	ErrInternalError = errors.New("an internal error has occurred")
)
