package crypto

import (
	"errors"
)

var (
	// ErrUnknownCrypto is returned in the event that the New()
	// function is called with the name of an implementation that
	// does not exist.
	ErrUnknownCrypto = errors.New("The specified crypto engine does not exist")

	// ErrInternalError is used to mask errors from the internal
	// crypto system that are unrecoverable.  This error is safe
	// to return whereas an error from a module may expose secure
	// information.
	ErrInternalError = errors.New("The crypto system has encountered an internal error")

	// ErrAuthorizationFailure is returned in the event the crypto
	// module determines that the provided secret does not match
	// the one secured earlier.
	ErrAuthorizationFailure = errors.New("Authorization failed - bad credentials")
)
