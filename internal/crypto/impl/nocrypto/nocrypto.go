package nocrypto

// THIS PACKAGE CONTAINS NO SECURITY CODE WHATSOEVER, DO NOT COMPILE
// IT INTO YOUR LIVE INSTALL!

import (
	"github.com/NetAuth/NetAuth/internal/crypto"
)

// NoCrypto binds the functions required by the crypto.Engine
// interface.
type NoCrypto struct{}

func init() {
	crypto.Register("nocrypto", New)
}

// New registers this crypto type for use by the NetAuth server.
func New() crypto.EMCrypto {
	return &NoCrypto{}
}

// SecureSecret returns the secret unmodified.
func (n *NoCrypto) SecureSecret(s string) (string, error) {
	return s, nil
}

// VerifySecret performs a string equality check to determine if the
// secret is legitimate.
func (n *NoCrypto) VerifySecret(s, h string) error {
	if s != h {
		return crypto.ErrAuthorizationFailure
	}
	return nil
}
