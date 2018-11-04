package hooks

import (
	"github.com/NetAuth/NetAuth/internal/crypto"

	pb "github.com/NetAuth/Protocol"
)

type SetEntitySecret struct {
	crypto.EMCrypto
}

func (*SetEntitySecret) Name() string  { return "set-entity-secret" }
func (*SetEntitySecret) Priority() int { return 50 }
func (s *SetEntitySecret) Run(e, de *pb.Entity) error {
	ssecret, err := s.SecureSecret(de.GetSecret())
	if err != nil {
		return err
	}
	e.Secret = &ssecret
	return nil
}
