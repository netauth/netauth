package crypto

import (
	"errors"
)

// The EMCrypto interface defines the functions that are needed to
// make a secret secure for storage and later verify a secret against
// the secured copy.
type EMCrypto interface {
	SecureSecret(string) (string, error)
	VerifySecret(string, string) error
}

// The Factory type is to be implemented by crypto implementations and
// shall be fed to the Register function.
type Factory func() EMCrypto

var (
	backends = make(map[string]Factory)

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

// New returns an initialized Crypto instance which can create and
// verify secure versions of secrets.
func New(name string) (EMCrypto, error) {
	b, ok := backends[name]
	if !ok {
		return nil, ErrUnknownCrypto
	}
	return b(), nil
}

// Register takes in a name for the engine and a function
// signature to bind to that name.
func Register(name string, newFunc Factory) {
	if _, ok := backends[name]; ok {
		// Return if the backend was already registered.
		return
	}
	backends[name] = newFunc
}

// GetBackendList returns a string list of hte backends that are
// available.
func GetBackendList() []string {
	var l []string

	for b := range backends {
		l = append(l, b)
	}

	return l
}
