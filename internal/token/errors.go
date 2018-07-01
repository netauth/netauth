package token

import (
	"errors"
)

var (
	// ErrUnknownTokenService is returned when a token name is
	// requested that isn't registered.
	ErrUnknownTokenService = errors.New("no token service with that name exists")

	// ErrKeyUnavailable signifies that at least one key is
	// unavailable to the token service.  For token systems that
	// use symmetric cryptography this is fatal, for token systems
	// that use asymmetric cryptography, this may be acceptable if
	// all you want to do is verify a token with a public key.
	ErrKeyUnavailable = errors.New("a required key is not available")

	// ErrKeyGenerationDisabled is returned when no keys were
	// available to load, and the option to generate keys has been
	// set false.
	ErrKeyGenerationDisabled = errors.New("key generation is disabled")

	// ErrInternalError captures all unidentified error cases
	// within various token services.
	ErrInternalError = errors.New("an unrecoverable internal error has occured")

	// ErrTokenInvalid is returned for generic cases where the
	// token is invalid for some reason.
	ErrTokenInvalid = errors.New("the provided token is invalid")
)
