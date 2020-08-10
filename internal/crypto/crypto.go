// Package crypto implements the plugin system for cryptography
// engines.  Specifically, the package implements methods to register
// cryptosystems and then obtain an initialized engine.
package crypto

import (
	"github.com/hashicorp/go-hclog"
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
type Factory func(hclog.Logger) (EMCrypto, error)

var (
	lb       hclog.Logger
	backends map[string]Factory
)

func init() {
	backends = make(map[string]Factory)
}

// New returns an initialized Crypto instance which can create and
// verify secure versions of secrets.
func New(backend string) (EMCrypto, error) {
	b, ok := backends[backend]
	if !ok {
		log().Error("Requested backend is not registered", "backend", backend)
		return nil, ErrUnknownCrypto
	}
	log().Info("Initializing backend", "backend", backend)
	return b(log())
}

// Register takes in a name for the engine and a function
// signature to bind to that name.
func Register(name string, newFunc Factory) {
	if _, ok := backends[name]; ok {
		// Return if the backend was already registered.
		log().Warn("A backend attempted to register an existing name", "backend", name)
		return
	}
	backends[name] = newFunc
	log().Info("Registered Backend", "backend", name)
}

// SetParentLogger sets the parent logger for this instance.
func SetParentLogger(l hclog.Logger) {
	lb = l.Named("crypto")
}

// log is a convenience function that will return a null logger if a
// parent logger has not been specified, mostly useful for tests.
func log() hclog.Logger {
	if lb == nil {
		lb = hclog.NewNullLogger()
	}
	return lb
}
