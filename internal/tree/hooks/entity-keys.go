package hooks

import (
	"context"

	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"
	"github.com/netauth/netauth/internal/tree/util"

	pb "github.com/netauth/protocol"
)

// ManageEntityKeys is a configurable hook that adds and removes keys.
type ManageEntityKeys struct {
	tree.BaseHook
	mode bool
}

// Run iterates on all keys in the request and adds or removes them
// from the entity's keystore.
func (mek *ManageEntityKeys) Run(_ context.Context, e, de *pb.Entity) error {
	for _, k := range de.Meta.Keys {
		e.Meta.Keys = util.PatchStringSlice(e.Meta.Keys, k, mek.mode, false)
	}
	return nil
}

func init() {
	startup.RegisterCallback(entityKeysCB)
}

func entityKeysCB() {
	tree.RegisterEntityHookConstructor("add-entity-key", NewAddEntityKey)
	tree.RegisterEntityHookConstructor("del-entity-key", NewDelEntityKey)
}

// NewAddEntityKey returns a hook initialized for adding keys.
func NewAddEntityKey(opts ...tree.HookOption) (tree.EntityHook, error) {
	opts = append([]tree.HookOption{
		tree.WithHookName("add-entity-key"),
		tree.WithHookPriority(50),
	}, opts...)

	return &ManageEntityKeys{tree.NewBaseHook(opts...), true}, nil
}

// NewDelEntityKey returns a hook initialized for removing keys.
func NewDelEntityKey(opts ...tree.HookOption) (tree.EntityHook, error) {
	opts = append([]tree.HookOption{
		tree.WithHookName("del-entity-key"),
		tree.WithHookPriority(50),
	}, opts...)

	return &ManageEntityKeys{tree.NewBaseHook(opts...), false}, nil
}
