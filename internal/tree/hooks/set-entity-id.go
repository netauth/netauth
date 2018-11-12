package hooks

import (
	"github.com/NetAuth/NetAuth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

// SetEntityID copies the ID from one entity to another.
type SetEntityID struct {
	tree.BaseHook
}

// Run copies the ID from de to e, no checks are enforced during the
// copy.
func (*SetEntityID) Run(e, de *pb.Entity) error {
	e.ID = de.ID
	return nil
}

func init() {
	tree.RegisterEntityHookConstructor("set-entity-id", NewSetEntityID)
}

// NewSetEntityID returns a SetEntityID hook initialized and ready for
// use.
func NewSetEntityID(c tree.RefContext) (tree.EntityProcessorHook, error) {
	return &SetEntityID{tree.NewBaseHook("set-entity-id", 50)}, nil
}
