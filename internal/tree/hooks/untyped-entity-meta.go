package hooks

import (
	"github.com/NetAuth/NetAuth/internal/tree"
	"github.com/NetAuth/NetAuth/internal/tree/util"

	pb "github.com/NetAuth/Protocol"
)

// ManageEntityUM is a configurable plugin that manages the untyped
// metadata for entities.
type ManageEntityUM struct {
	tree.BaseHook
	mode string
}

// Run will process metadata provided via de onto e according to the
// mode the plugin is configured for.  "UPSERT" will add or update
// fields as appropriate.  "CLEARFUZZY" will ignore Z-Indexing
// annotations.  "CLEAREXACT" will require exact key specifications.
func (mm *ManageEntityUM) Run(e, de *pb.Entity) error {
	for _, m := range de.Meta.UntypedMeta {
		key, value := splitKeyValue(m)
		e.Meta.UntypedMeta = util.PatchKeyValueSlice(e.Meta.UntypedMeta, mm.mode, key, value)
	}
	return nil
}

func init() {
	tree.RegisterEntityHookConstructor("add-untyped-metadata", NewAddEntityUM)
	tree.RegisterEntityHookConstructor("del-untyped-metadata-fuzzy", NewDelFuzzyEntityUM)
	tree.RegisterEntityHookConstructor("del-untyped-metadata-exact", NewDelExactEntityUM)
}

// NewAddEntityUM returns a configured hook in UPSERT mode.
func NewAddEntityUM(c tree.RefContext) (tree.EntityHook, error) {
	return &ManageEntityUM{tree.NewBaseHook("add-untyped-metadata", 50), "UPSERT"}, nil
}

// NewDelFuzzyEntityUM returns a configured hook in CLEARFUZZY mode.
func NewDelFuzzyEntityUM(c tree.RefContext) (tree.EntityHook, error) {
	return &ManageEntityUM{tree.NewBaseHook("del-untyped-metadata-fuzzy", 50), "CLEARFUZZY"}, nil
}

// NewDelExactEntityUM returns a configured hook in CLEAREXACT mode.
func NewDelExactEntityUM(c tree.RefContext) (tree.EntityHook, error) {
	return &ManageEntityUM{tree.NewBaseHook("del-untyped-metadata-exact", 50), "CLEAREXACT"}, nil
}
