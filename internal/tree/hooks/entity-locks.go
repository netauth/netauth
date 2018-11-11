package hooks

import (
	"github.com/NetAuth/NetAuth/internal/tree"
	"github.com/golang/protobuf/proto"

	pb "github.com/NetAuth/Protocol"
)

type EntityLockManager struct {
	tree.BaseHook
	lockstate bool
}

func (elm *EntityLockManager) Run(e, de *pb.Entity) error {
	e.Meta.Locked = proto.Bool(elm.lockstate)
	return nil
}

func init() {
	tree.RegisterEntityHookConstructor("lock-entity", NewELMLock)
	tree.RegisterEntityHookConstructor("unlock-entity", NewELMUnlock)
}

func NewELMLock(c tree.RefContext) (tree.EntityProcessorHook, error) {
	return &EntityLockManager{tree.NewBaseHook("lock-entity", 40), true}, nil
}

func NewELMUnlock(c tree.RefContext) (tree.EntityProcessorHook, error) {
	return &EntityLockManager{tree.NewBaseHook("unlock-entity", 40), false}, nil
}
