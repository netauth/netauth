package bcrypt

import (
	"flag"
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/NetAuth/NetAuth/internal/server/crypto"
	"github.com/NetAuth/NetAuth/pkg/errors"
)

var (
	cost = flag.Int("bcrypt_cost", 15, "Cost to use when running the bcrypt hashing algorithm")
)

func init() {
	crypto.RegisterCrypto("bcrypt", New)
}

type BCryptEngine struct {
	cost int
}

func New() crypto.EMCrypto {
	x := new(BCryptEngine)
	x.cost = *cost
	return x
}

func (b *BCryptEngine) SecureSecret(secret string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(secret), b.cost)
	if err != nil {
		log.Printf("Crypto Fault: %s", err)
		return "", errors.E_CRYPTO_FAULT
	}
	return string(hash[:]), nil
}

func (b *BCryptEngine) VerifySecret(secret, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(secret))
	if err != nil {
		log.Printf("Crypto Error: %s", err)
		return errors.E_CRYPTO_BADAUTH
	}
	return nil
}
