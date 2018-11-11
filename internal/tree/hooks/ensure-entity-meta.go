package hooks

import (
	"github.com/NetAuth/NetAuth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

type EnsureEntityMeta struct {
	tree.BaseHook
}

func (*EnsureEntityMeta) Run(e, de *pb.Entity) error {
	if e.Meta == nil {
		e.Meta = &pb.EntityMeta{}
	}
	return nil
}

func init() {
	tree.RegisterEntityHookConstructor("ensure-entity-meta", NewEnsureEntityMeta)
}

func NewEnsureEntityMeta(c tree.RefContext) (tree.EntityProcessorHook, error) {
	return &EnsureEntityMeta{tree.NewBaseHook("ensure-entity-meta", 20)}, nil
}
