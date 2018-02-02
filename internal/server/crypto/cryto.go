package crypto

import (
	"github.com/NetAuth/NetAuth/pkg/errors"
)

type EMCrypto interface {
	SecureSecret(string) (string, error)
	VerifySecret(string, string) error
}

type CryptoFactory func() EMCrypto

var (
	backends = make(map[string]CryptoFactory)
)

// New returns an initialized Crypto instance which can create and
// verify secure versions of secrets.
func New(name string) (EMCrypto, error) {
	b, ok := backends[name]
	if !ok {
		return nil, errors.E_NO_SUCH_CRYPTO
	}
	return b(), nil
}

// RegisterCrypto takes in a name for the engine and a function
// signature to bind to that name.
func RegisterCrypto(name string, newFunc CryptoFactory) {
	if _, ok := backends[name]; ok {
		// Return if the backend was already registered.
		return
	}
	backends[name] = newFunc
}

// GetBackendList returns a string list of hte backends that are
// available.
func GetBackendList() []string {
	l := []string{}

	for b, _ := range backends {
		l = append(l, b)
	}

	return l
}
