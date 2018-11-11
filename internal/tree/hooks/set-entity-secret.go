package hooks

import (
	"github.com/NetAuth/NetAuth/internal/crypto"
	"github.com/NetAuth/NetAuth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

type SetEntitySecret struct {
	tree.BaseHook
	crypto.EMCrypto
}

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

func NewSetEntitySecret(c tree.RefContext) (tree.EntityProcessorHook, error) {
	return &SetEntitySecret{tree.NewBaseHook("set-entity-secret", 50), c.Crypto}, nil
}
