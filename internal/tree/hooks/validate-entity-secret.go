package hooks

import (
	"context"

	"github.com/netauth/netauth/internal/crypto"
	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

// ValidateEntitySecret passes the secret to the crypto engine for
// validation.
type ValidateEntitySecret struct {
	tree.BaseHook
	crypto.EMCrypto
}

// Run calls VerifySecret to compare de.Secret with the secured copy from e.Secret.
func (v *ValidateEntitySecret) Run(_ context.Context, e, de *pb.Entity) error {
	return v.VerifySecret(de.GetSecret(), e.GetSecret())
}

func init() {
	startup.RegisterCallback(validateEntitySecretCB)
}

func validateEntitySecretCB() {
	tree.RegisterEntityHookConstructor("validate-entity-secret", NewValidateEntitySecret)
}

// NewValidateEntitySecret returns an initialized hook ready for use.
func NewValidateEntitySecret(c tree.RefContext) (tree.EntityHook, error) {
	return &ValidateEntitySecret{tree.NewBaseHook("validate-entity-secret", 50), c.Crypto}, nil
}
