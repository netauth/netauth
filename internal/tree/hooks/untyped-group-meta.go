package hooks

import (
	"github.com/NetAuth/NetAuth/internal/tree"
	"github.com/NetAuth/NetAuth/internal/tree/util"

	pb "github.com/NetAuth/Protocol"
)

// ManageGroupUM is a configurable plugin that manages the untyped
// metadata for groups.
type ManageGroupUM struct {
	tree.BaseHook
	mode string
}

// Run will process metadata provided via dg onto g according to the
// mode the plugin is configured for.  "UPSERT" will add or update
// fields as appropriate.  "CLEARFUZZY" will ignore Z-Indexing
// annotations.  "CLEAREXACT" will require exact key specifications.
func (mm *ManageGroupUM) Run(g, dg *pb.Group) error {
	for _, m := range dg.UntypedMeta {
		key, value := splitKeyValue(m)
		g.UntypedMeta = util.PatchKeyValueSlice(g.UntypedMeta, mm.mode, key, value)
	}
	return nil
}

func init() {
	tree.RegisterGroupHookConstructor("add-untyped-metadata", NewAddGroupUM)
	tree.RegisterGroupHookConstructor("del-untyped-metadata-fuzzy", NewDelFuzzyGroupUM)
	tree.RegisterGroupHookConstructor("del-untyped-metadata-exact", NewDelExactGroupUM)
}

// NewAddGroupUM returns a configured hook in UPSERT mode.
func NewAddGroupUM(c tree.RefContext) (tree.GroupHook, error) {
	return &ManageGroupUM{tree.NewBaseHook("add-untyped-metadata", 50), "UPSERT"}, nil
}

// NewDelFuzzyGroupUM returns a configured hook in CLEARFUZZY mode.
func NewDelFuzzyGroupUM(c tree.RefContext) (tree.GroupHook, error) {
	return &ManageGroupUM{tree.NewBaseHook("del-untyped-metadata-fuzzy", 50), "CLEARFUZZY"}, nil
}

// NewDelExactGroupUM returns a configured hook in CLEAREXACT mode.
func NewDelExactGroupUM(c tree.RefContext) (tree.GroupHook, error) {
	return &ManageGroupUM{tree.NewBaseHook("del-untyped-metadata-exact", 50), "CLEAREXACT"}, nil
}
