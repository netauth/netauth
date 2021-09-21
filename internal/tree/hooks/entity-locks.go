package hooks

import (
	"context"

	"google.golang.org/protobuf/proto"

	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

// The EntityLockManager is a configurable hook that can either lock
// or unlock entities as needed.
type EntityLockManager struct {
	tree.BaseHook
	lockstate bool
}

// Run will set the entity lock status unconditionally to the
// configured value for the instantiated hook.
func (elm *EntityLockManager) Run(_ context.Context, e, de *pb.Entity) error {
	e.Meta.Locked = proto.Bool(elm.lockstate)
	return nil
}

func init() {
	startup.RegisterCallback(entityLockCB)
}

func entityLockCB() {
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
