package bcrypt

import (
	"github.com/hashicorp/go-hclog"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"

	"github.com/netauth/netauth/internal/crypto"
)

func init() {
	crypto.RegisterCallback(cb)
	pflag.Int("crypto.bcrypt.cost", 15, "Cost for bcrypt")
}

func cb() {
	crypto.Register("bcrypt", New)
}

// Engine binds the functions of the BCrypt Crypto system and
// satisfies the crypto.EMCrypto interface.  The parameter 'cost' is
// used to set the cost of the algorithm and should be set on a
// per-site basis.
type Engine struct {
	cost int
	l    hclog.Logger
}

// New registers this crypto type for use by the NetAuth server.
func New(l hclog.Logger) (crypto.EMCrypto, error) {
	x := new(Engine)
	x.cost = viper.GetInt("crypto.bcrypt.cost")
	x.l = l.Named("bcrypt")
	x.l.Debug("BCrypt Initialized", "cost", x.cost)
	return x, nil
}

// SecureSecret takes in a secret and generates a bcrypt hash from it.
// This is then returned for storage in the database.
func (b *Engine) SecureSecret(secret string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(secret), b.cost)
	if err != nil {
		b.l.Error("Bcrypt Error has occurred", "error", err)
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
		b.l.Error("Bcrypt Error has occurred", "error", err)
		return crypto.ErrAuthorizationFailure
	}
	return nil
}
