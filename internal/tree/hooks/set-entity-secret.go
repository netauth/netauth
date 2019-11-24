package hooks

import (
	"github.com/netauth/netauth/internal/crypto"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

// SetEntitySecret takes a plaintext secret and converts it to a
// secured secret for storage.
type SetEntitySecret struct {
	tree.BaseHook
	crypto.EMCrypto
}

// Run takes a plaintext secret from de.Secret and secures it using a
// crypto engine.  The secured secret will be written to e.Secret.
func (s *SetEntitySecret) Run(e, de *pb.Entity) error {
	ssecret, err := s.SecureSecret(de.GetSecret())
	if err != nil {
		return err
	}
	e.Secret = &ssecret
	return nil
}

func init() {
	tree.RegisterEntityHookConstructor("set-entity-secret", NewSetEntitySecret)
}

// NewSetEntitySecret returns an initialized hook for use.
func NewSetEntitySecret(c tree.RefContext) (tree.EntityHook, error) {
	return &SetEntitySecret{tree.NewBaseHook("set-entity-secret", 50), c.Crypto}, nil
}
