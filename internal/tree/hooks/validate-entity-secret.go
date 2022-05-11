package hooks

import (
	"context"

	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

// ValidateEntitySecret passes the secret to the crypto engine for
// validation.
type ValidateEntitySecret struct {
	tree.BaseHook
}

// Run calls VerifySecret to compare de.Secret with the secured copy from e.Secret.
func (v *ValidateEntitySecret) Run(_ context.Context, e, de *pb.Entity) error {
	return v.Crypto().VerifySecret(de.GetSecret(), e.GetSecret())
}

func init() {
	startup.RegisterCallback(validateEntitySecretCB)
}

func validateEntitySecretCB() {
	tree.RegisterEntityHookConstructor("validate-entity-secret", NewValidateEntitySecret)
}

// NewValidateEntitySecret returns an initialized hook ready for use.
func NewValidateEntitySecret(opts ...tree.HookOption) (tree.EntityHook, error) {
	opts = append([]tree.HookOption{
		tree.WithHookName("validate-entity-secret"),
		tree.WithHookPriority(50),
	}, opts...)

	return &ValidateEntitySecret{tree.NewBaseHook(opts...)}, nil
}
