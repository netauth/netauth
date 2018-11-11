package hooks

import (
	"github.com/NetAuth/NetAuth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

type ValidateEntityUnlocked struct {
	tree.BaseHook
}

func (*ValidateEntityUnlocked) Run(e, de *pb.Entity) error {
	if e.GetMeta().GetLocked() {
		return tree.ErrEntityLocked
	}
	return nil
}

func init() {
	tree.RegisterEntityHookConstructor("validate-entity-unlocked", NewValidateEntityUnlocked)
}

func NewValidateEntityUnlocked(c tree.RefContext) (tree.EntityProcessorHook, error) {
	return &ValidateEntityUnlocked{tree.NewBaseHook("validate-entity-unlocked", 20)}, nil
}
