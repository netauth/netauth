package bcrypt

import (
	"flag"
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/NetAuth/NetAuth/internal/crypto"
)

var (
	cost = flag.Int("bcrypt_cost", 15, "Cost to use when running the bcrypt hashing algorithm")
)

func init() {
	crypto.Register("bcrypt", New)
}

// Engine binds the functions of the BCrypt Crypto system and
// satisfies the crypto.EMCrypto interface.  The parameter 'cost' is
// used to set the cost of the algorithm and should be set on a
// per-site basis.
type Engine struct {
	cost int
}

// New registers this crypto type for use by the NetAuth server.
func New() crypto.EMCrypto {
	x := new(Engine)
	x.cost = *cost
	return x
}

// SecureSecret takes in a secret and generates a bcrypt hash from it.
// This is then returned for storage in the database.
func (b *Engine) SecureSecret(secret string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(secret), b.cost)
	if err != nil {
		log.Printf("Crypto Fault: %s", err)
		return "", crypto.ErrInternalError
	}
	return string(hash[:]), nil
}

// VerifySecret verifies a given secret against a given hash and
// returns either nil for a match or a crypto.ErrAuthorizationFailure
// in the case that the secret did not match the stored one.
func (b *Engine) VerifySecret(secret, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(secret))
	if err != nil {
		log.Printf("Crypto Error: %s", err)
		return crypto.ErrAuthorizationFailure
	}
	return nil
}
