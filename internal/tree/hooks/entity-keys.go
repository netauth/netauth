package hooks

import (
	"github.com/NetAuth/NetAuth/internal/tree"
	"github.com/NetAuth/NetAuth/internal/tree/util"

	pb "github.com/NetAuth/Protocol"
)

type ManageEntityKeys struct {
	tree.BaseHook
	mode bool
}

func (mek *ManageEntityKeys) Run(e, de *pb.Entity) error {
	for _, k := range de.Meta.Keys {
		e.Meta.Keys = util.PatchStringSlice(e.Meta.Keys, k, true, mek.mode)
	}
	return nil
}

func init() {
	tree.RegisterEntityHookConstructor("add-entity-key", NewAddEntityKey)
	tree.RegisterEntityHookConstructor("del-entity-key", NewDelEntityKey)
}

func NewAddEntityKey(c tree.RefContext) (tree.EntityProcessorHook, error) {
	return &ManageEntityKeys{tree.NewBaseHook("add-entity-key", 50), true}, nil
}

func NewDelEntityKey(c tree.RefContext) (tree.EntityProcessorHook, error) {
	return &ManageEntityKeys{tree.NewBaseHook("del-entity-key", 50), false}, nil
}
