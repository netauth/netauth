package hooks

import (
	"context"

	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

// SetEntitySecret takes a plaintext secret and converts it to a
// secured secret for storage.
type SetEntitySecret struct {
	tree.BaseHook
}

// Run takes a plaintext secret from de.Secret and secures it using a
// crypto engine.  The secured secret will be written to e.Secret.
func (s *SetEntitySecret) Run(_ context.Context, e, de *pb.Entity) error {
	ssecret, err := s.Crypto().SecureSecret(de.GetSecret())
	if err != nil {
		return err
	}
	e.Secret = &ssecret
	return nil
}

func init() {
	startup.RegisterCallback(setEntitySecretCB)
}

func setEntitySecretCB() {
	tree.RegisterEntityHookConstructor("set-entity-secret", NewSetEntitySecret)
}

// NewSetEntitySecret returns an initialized hook for use.
func NewSetEntitySecret(opts ...tree.HookOption) (tree.EntityHook, error) {
	opts = append([]tree.HookOption{
		tree.WithHookName("set-entity-secret"),
		tree.WithHookPriority(50),
	}, opts...)

	return &SetEntitySecret{tree.NewBaseHook(opts...)}, nil
}
