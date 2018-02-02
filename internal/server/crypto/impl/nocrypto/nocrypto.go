package nocrypto

// THIS PACKAGE CONTAINS NO SECURITY CODE WHATSOEVER, DO NOT COMPILE
// IT INTO YOUR LIVE INSTALL!

import (
	"github.com/NetAuth/NetAuth/internal/server/crypto"
	"github.com/NetAuth/NetAuth/pkg/errors"
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
		return errors.E_CRYPTO_BADAUTH
	}
	return nil
}
