package hooks

import (
	"github.com/NetAuth/NetAuth/internal/crypto"
	"github.com/NetAuth/NetAuth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

type ValidateEntitySecret struct {
	tree.BaseHook
	crypto.EMCrypto
}

func (v *ValidateEntitySecret) Run(e, de *pb.Entity) error {
	if err := v.VerifySecret(de.GetSecret(), e.GetSecret()); err != nil {
		return err
	}
	return nil
}

func init() {
	tree.RegisterEntityHookConstructor("validate-entity-secret", NewValidateEntitySecret)
}

func NewValidateEntitySecret(c tree.RefContext) (tree.EntityProcessorHook, error) {
	return &ValidateEntitySecret{tree.NewBaseHook("validate-entity-secret", 50), c.Crypto}, nil
}
