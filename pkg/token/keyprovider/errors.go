package keyprovider

import "errors"

var (
	// ErrUnknownKeyProvider is returned when a key provider is
	// requested for which no corresponding factory has been
	// registered.  Check your initialization order if this is
	// returned unexpectedly.
	ErrUnknownKeyProvider = errors.New("no key provider with that name")

	// ErrNoSuchKey is returned when the key requested can't be
	// retrieved.
	ErrNoSuchKey = errors.New("no key exists with that identifier")

	// ErrInternal is returned when something screwy is happening
	// and a key can't be provided for some non-recoverable error.
	ErrInternal = errors.New("an internal error has occured")
)
