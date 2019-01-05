// Package crypto implements the plugin system for cryptography
// engines.  Specifically, the package implements methods to register
// cryptosystems and then obtain an initialized engine.
package crypto

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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
type Factory func() (EMCrypto, error)

var (
	backends map[string]Factory
)

func init() {
	backends = make(map[string]Factory)
	pflag.String("crypto.backend", "bcrypt", "Cryptography system to use")
}

// New returns an initialized Crypto instance which can create and
// verify secure versions of secrets.
func New() (EMCrypto, error) {
	b, ok := backends[viper.GetString("crypto.backend")]
	if !ok {
		return nil, ErrUnknownCrypto
	}
	return b()
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
