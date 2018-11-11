package hooks

import (
	"strings"

	"github.com/NetAuth/NetAuth/internal/tree"
	"github.com/NetAuth/NetAuth/internal/tree/util"

	pb "github.com/NetAuth/Protocol"
)

type ManageEntityUM struct {
	tree.BaseHook
	mode string
}

func (mm *ManageEntityUM) Run(e, de *pb.Entity) error {
	for _, m := range de.Meta.UntypedMeta {
		key, value := splitKeyValue(m)
		e.Meta.UntypedMeta = util.PatchKeyValueSlice(e.Meta.UntypedMeta, mm.mode, key, value)
	}
	return nil
}

func splitKeyValue(s string) (string, string) {
	parts := strings.SplitN(s, ":", 2)
	return parts[0], parts[1]
}

func init() {
	tree.RegisterEntityHookConstructor("add-untyped-metadata", NewAddEntityUM)
	tree.RegisterEntityHookConstructor("del-untyped-metadata-fuzzy", NewDelFuzzyEntityUM)
	tree.RegisterEntityHookConstructor("del-untyped-metadata-exact", NewDelExactEntityUM)
}

func NewAddEntityUM(c tree.RefContext) (tree.EntityProcessorHook, error) {
	return &ManageEntityUM{tree.NewBaseHook("add-untyped-metadata", 50), "UPSERT"}, nil
}

func NewDelFuzzyEntityUM(c tree.RefContext) (tree.EntityProcessorHook, error) {
	return &ManageEntityUM{tree.NewBaseHook("del-untyped-metadata-fuzzy", 50), "CLEARFUZZY"}, nil
}

func NewDelExactEntityUM(c tree.RefContext) (tree.EntityProcessorHook, error) {
	return &ManageEntityUM{tree.NewBaseHook("del-untyped-metadata-exact", 50), "UPSERT"}, nil
}
