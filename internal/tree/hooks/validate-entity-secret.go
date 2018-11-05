package hooks

import (
	"github.com/NetAuth/NetAuth/internal/crypto"

	pb "github.com/NetAuth/Protocol"
)

type ValidateEntitySecret struct {
	crypto.EMCrypto
}

func (*ValidateEntitySecret) Name() string  { return "validate-entity-secret" }
func (*ValidateEntitySecret) Priority() int { return 40 }
func (v *ValidateEntitySecret) Run(e, de *pb.Entity) error {
	if err := v.VerifySecret(de.GetSecret(), e.GetSecret()); err != nil {
		return err
	}
	return nil
}
