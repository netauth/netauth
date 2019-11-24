package hooks

import (
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
func (mek *ManageEntityKeys) Run(e, de *pb.Entity) error {
	for _, k := range de.Meta.Keys {
		e.Meta.Keys = util.PatchStringSlice(e.Meta.Keys, k, mek.mode, false)
	}
	return nil
}

func init() {
	tree.RegisterEntityHookConstructor("add-entity-key", NewAddEntityKey)
	tree.RegisterEntityHookConstructor("del-entity-key", NewDelEntityKey)
}

// NewAddEntityKey returns a hook initialized for adding keys.
func NewAddEntityKey(c tree.RefContext) (tree.EntityHook, error) {
	return &ManageEntityKeys{tree.NewBaseHook("add-entity-key", 50), true}, nil
}

// NewDelEntityKey returns a hook initialized for removing keys.
func NewDelEntityKey(c tree.RefContext) (tree.EntityHook, error) {
	return &ManageEntityKeys{tree.NewBaseHook("del-entity-key", 50), false}, nil
}
