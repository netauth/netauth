package hooks

import (
	"context"

	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

// DestroyGroup removes an entity from the system.
type DestroyGroup struct {
	tree.BaseHook
}

// Run will request the underlying datastore to remove the group,
// returning any status provided.  If the group Name is not specified
// in g, it will be obtained from dg.
func (d *DestroyGroup) Run(ctx context.Context, g, dg *pb.Group) error {
	// This hook is somewhat special since it might be called
	// after a processing pipeline, or just to remove a group.
	if g.GetName() == "" {
		g.Name = dg.Name
	}
	return d.Storage().DeleteGroup(ctx, g.GetName())
}

func init() {
	startup.RegisterCallback(destroyGroupCB)
}

func destroyGroupCB() {
	tree.RegisterGroupHookConstructor("destroy-group", NewDestroyGroup)
}

// NewDestroyGroup returns an initialized DestroyGroup hook for use.
func NewDestroyGroup(opts ...tree.HookOption) (tree.GroupHook, error) {
	opts = append([]tree.HookOption{
		tree.WithHookName("destroy-group"),
		tree.WithHookPriority(99),
	}, opts...)
	return &DestroyGroup{tree.NewBaseHook(opts...)}, nil
}
