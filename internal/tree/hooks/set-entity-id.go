package hooks

import (
	"github.com/NetAuth/NetAuth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

type SetEntityID struct {
	tree.BaseHook
}

func (*SetEntityID) Run(e, de *pb.Entity) error {
	e.ID = de.ID
	return nil
}

func init() {
	tree.RegisterEntityHookConstructor("set-entity-id", NewSetEntityID)
}

func NewSetEntityID(c tree.RefContext) (tree.EntityProcessorHook, error) {
	return &SetEntityID{tree.NewBaseHook("set-entity-id", 50)}, nil
}
