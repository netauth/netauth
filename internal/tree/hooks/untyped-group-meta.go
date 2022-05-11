package hooks

import (
	"context"

	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"
	"github.com/netauth/netauth/internal/tree/util"

	pb "github.com/netauth/protocol"
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
func (mm *ManageGroupUM) Run(_ context.Context, g, dg *pb.Group) error {
	for _, m := range dg.UntypedMeta {
		key, value := splitKeyValue(m)
		g.UntypedMeta = util.PatchKeyValueSlice(g.UntypedMeta, mm.mode, key, value)
	}
	return nil
}

func init() {
	startup.RegisterCallback(manageGroupUMCB)
}

func manageGroupUMCB() {
	tree.RegisterGroupHookConstructor("add-untyped-metadata", NewAddGroupUM)
	tree.RegisterGroupHookConstructor("del-untyped-metadata-fuzzy", NewDelFuzzyGroupUM)
	tree.RegisterGroupHookConstructor("del-untyped-metadata-exact", NewDelExactGroupUM)
}

// NewAddGroupUM returns a configured hook in UPSERT mode.
func NewAddGroupUM(opts ...tree.HookOption) (tree.GroupHook, error) {
	opts = append([]tree.HookOption{
		tree.WithHookName("add-untyped-metadata"),
		tree.WithHookPriority(50),
	}, opts...)

	return &ManageGroupUM{tree.NewBaseHook(opts...), "UPSERT"}, nil
}

// NewDelFuzzyGroupUM returns a configured hook in CLEARFUZZY mode.
func NewDelFuzzyGroupUM(opts ...tree.HookOption) (tree.GroupHook, error) {
	opts = append([]tree.HookOption{
		tree.WithHookName("del-untyped-metadata-fuzzy"),
		tree.WithHookPriority(50),
	}, opts...)

	return &ManageGroupUM{tree.NewBaseHook(opts...), "CLEARFUZZY"}, nil
}

// NewDelExactGroupUM returns a configured hook in CLEAREXACT mode.
func NewDelExactGroupUM(opts ...tree.HookOption) (tree.GroupHook, error) {
	opts = append([]tree.HookOption{
		tree.WithHookName("del-untyped-metadata-exact"),
		tree.WithHookPriority(50),
	}, opts...)

	return &ManageGroupUM{tree.NewBaseHook(opts...), "CLEAREXACT"}, nil
}
