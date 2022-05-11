package hooks

import (
	"context"

	"github.com/netauth/netauth/internal/tree"
	"github.com/netauth/netauth/internal/tree/util"

	pb "github.com/netauth/protocol"
)

// DirectGroupManager implements a hook type that can add or remove
// the groups that an entity is directly a member of.
type DirectGroupManager struct {
	tree.BaseHook
	mode bool
}

// Run iterates through all groups in de.Meta.Groups and adds or
// removes them from e.Meta.Groups based on the value of dgm.mode.
// True will add groups, false will remove them.
func (dgm *DirectGroupManager) Run(_ context.Context, e, de *pb.Entity) error {
	groups := de.GetMeta().GetGroups()
	for i := range groups {
		// Patch the group list and match groups exactly.
		e.Meta.Groups = util.PatchStringSlice(e.Meta.Groups, groups[i], dgm.mode, true)
	}
	return nil
}

func init() {
	tree.RegisterEntityHookConstructor("add-direct-group", NewAddDirectGroup)
	tree.RegisterEntityHookConstructor("del-direct-group", NewDelDirectGroup)
}

// NewAddDirectGroup returns a DirectGroupManager initialized in add
// mode.
func NewAddDirectGroup(opts ...tree.HookOption) (tree.EntityHook, error) {
	opts = append([]tree.HookOption{
		tree.WithHookName("add-direct-group"),
		tree.WithHookPriority(50),
	}, opts...)
	return &DirectGroupManager{tree.NewBaseHook(opts...), true}, nil
}

// NewDelDirectGroup returns a DirectGroupManager initialized in
// delete mode.
func NewDelDirectGroup(opts ...tree.HookOption) (tree.EntityHook, error) {
	opts = append([]tree.HookOption{
		tree.WithHookName("del-direct-group"),
		tree.WithHookPriority(50),
	}, opts...)
	return &DirectGroupManager{tree.NewBaseHook(opts...), false}, nil
}
