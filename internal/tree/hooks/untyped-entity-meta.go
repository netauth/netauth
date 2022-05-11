package hooks

import (
	"context"

	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"
	"github.com/netauth/netauth/internal/tree/util"

	pb "github.com/netauth/protocol"
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
func (mm *ManageEntityUM) Run(_ context.Context, e, de *pb.Entity) error {
	for _, m := range de.Meta.UntypedMeta {
		key, value := splitKeyValue(m)
		e.Meta.UntypedMeta = util.PatchKeyValueSlice(e.Meta.UntypedMeta, mm.mode, key, value)
	}
	return nil
}

func init() {
	startup.RegisterCallback(manageEntityUMCB)
}

func manageEntityUMCB() {
	tree.RegisterEntityHookConstructor("add-untyped-metadata", NewAddEntityUM)
	tree.RegisterEntityHookConstructor("del-untyped-metadata-fuzzy", NewDelFuzzyEntityUM)
	tree.RegisterEntityHookConstructor("del-untyped-metadata-exact", NewDelExactEntityUM)
}

// NewAddEntityUM returns a configured hook in UPSERT mode.
func NewAddEntityUM(opts ...tree.HookOption) (tree.EntityHook, error) {
	opts = append([]tree.HookOption{
		tree.WithHookName("add-untyped-metadata"),
		tree.WithHookPriority(50),
	}, opts...)

	return &ManageEntityUM{tree.NewBaseHook(opts...), "UPSERT"}, nil
}

// NewDelFuzzyEntityUM returns a configured hook in CLEARFUZZY mode.
func NewDelFuzzyEntityUM(opts ...tree.HookOption) (tree.EntityHook, error) {
	opts = append([]tree.HookOption{
		tree.WithHookName("del-untyped-metadata-fuzzy"),
		tree.WithHookPriority(50),
	}, opts...)

	return &ManageEntityUM{tree.NewBaseHook(opts...), "CLEARFUZZY"}, nil
}

// NewDelExactEntityUM returns a configured hook in CLEAREXACT mode.
func NewDelExactEntityUM(opts ...tree.HookOption) (tree.EntityHook, error) {
	opts = append([]tree.HookOption{
		tree.WithHookName("del-untyped-metadata-exact"),
		tree.WithHookPriority(50),
	}, opts...)

	return &ManageEntityUM{tree.NewBaseHook(opts...), "CLEAREXACT"}, nil
}
