package nocrypto

// THIS PACKAGE CONTAINS NO SECURITY CODE WHATSOEVER, DO NOT COMPILE
// IT INTO YOUR LIVE INSTALL!

import (
	"github.com/NetAuth/NetAuth/internal/crypto"
)

type NoCrypto struct{}

func init() {
	crypto.RegisterCrypto("nocrypto", New)
}

func New() crypto.EMCrypto {
	return & NoCrypto{}
}

func (n *NoCrypto) SecureSecret(s string) (string, error) {
	return s, nil
}

func (n *NoCrypto) VerifySecret(s, h string) error {
	if s != h {
		return crypto.AuthorizationFailure
	}
	return nil
}
