package hooks

import (
	"github.com/NetAuth/NetAuth/internal/tree/errors"

	pb "github.com/NetAuth/Protocol"
)

type ValidateEntityUnlocked struct{}

func (*ValidateEntityUnlocked) Name() string  { return "validate-entity-unlocked" }
func (*ValidateEntityUnlocked) Priority() int { return 25 }
func (*ValidateEntityUnlocked) Run(e, de *pb.Entity) error {
	if e.GetMeta().GetLocked() {
		return tree.ErrEntityLocked
	}
	return nil
}
