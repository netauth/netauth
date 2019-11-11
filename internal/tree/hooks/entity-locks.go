package hooks

import (
	"github.com/golang/protobuf/proto"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

// The EntityLockManager is a configurable hook that can either lock
// or unlock entities as needed.
type EntityLockManager struct {
	tree.BaseHook
	lockstate bool
}

// Run will set the entity lock status unconditionally to the
// configured value for the instantiated hook.
func (elm *EntityLockManager) Run(e, de *pb.Entity) error {
	e.Meta.Locked = proto.Bool(elm.lockstate)
	return nil
}

func init() {
	tree.RegisterEntityHookConstructor("lock-entity", NewELMLock)
	tree.RegisterEntityHookConstructor("unlock-entity", NewELMUnlock)
}

// NewELMLock returns a configured hook in LOCK mode.
func NewELMLock(c tree.RefContext) (tree.EntityHook, error) {
	return &EntityLockManager{tree.NewBaseHook("lock-entity", 40), true}, nil
}

// NewELMUnlock returns a configured hook in UNLOCK mode.
func NewELMUnlock(c tree.RefContext) (tree.EntityHook, error) {
	return &EntityLockManager{tree.NewBaseHook("unlock-entity", 40), false}, nil
}
