package hooks

import (
	"context"

	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

// ValidateEntityUnlocked returns an error if the entity is locked.
type ValidateEntityUnlocked struct {
	tree.BaseHook
}

// Run queries the locked status of an entity and returns either
// ErrEntityLocked or nil, depending on if the entity is locked.
func (*ValidateEntityUnlocked) Run(_ context.Context, e, de *pb.Entity) error {
	if e.GetMeta().GetLocked() {
		return tree.ErrEntityLocked
	}
	return nil
}

func init() {
	startup.RegisterCallback(validateEntityUnlockedCB)
}

func validateEntityUnlockedCB() {
	tree.RegisterEntityHookConstructor("validate-entity-unlocked", NewValidateEntityUnlocked)
}

// NewValidateEntityUnlocked returns an initialized hook.
func NewValidateEntityUnlocked(c tree.RefContext) (tree.EntityHook, error) {
	return &ValidateEntityUnlocked{tree.NewBaseHook("validate-entity-unlocked", 20)}, nil
}
