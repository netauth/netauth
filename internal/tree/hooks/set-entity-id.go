package hooks

import (
	"context"

	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

// SetEntityID copies the ID from one entity to another.
type SetEntityID struct {
	tree.BaseHook
}

// Run copies the ID from de to e, no checks are enforced during the
// copy.
func (*SetEntityID) Run(_ context.Context, e, de *pb.Entity) error {
	e.ID = de.ID
	return nil
}

func init() {
	startup.RegisterCallback(setEntityIDCB)
}

func setEntityIDCB() {
	tree.RegisterEntityHookConstructor("set-entity-id", NewSetEntityID)
}

// NewSetEntityID returns a SetEntityID hook initialized and ready for
// use.
func NewSetEntityID(c tree.RefContext) (tree.EntityHook, error) {
	return &SetEntityID{tree.NewBaseHook("set-entity-id", 50)}, nil
}
