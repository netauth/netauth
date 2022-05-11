package hooks

import (
	"context"

	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

// FailOnExistingGroup is a hook that can be used to guard creation
// processes on groups.
type FailOnExistingGroup struct {
	tree.BaseHook
}

// Run contacts the datastore and attempts to load the group specified
// by dg.  If the group loads successfully then an error is returned,
// in other cases nil is returned.
func (f *FailOnExistingGroup) Run(ctx context.Context, g, dg *pb.Group) error {
	if _, err := f.Storage().LoadGroup(ctx, dg.GetName()); err == nil {
		return tree.ErrDuplicateGroupName
	}
	return nil
}

func init() {
	startup.RegisterCallback(failOnExistingGroupCB)
}

func failOnExistingGroupCB() {
	tree.RegisterGroupHookConstructor("fail-on-existing-group", NewFailOnExistingGroup)
}

// NewFailOnExistingGroup returns an initialized hook ready for use.
func NewFailOnExistingGroup(opts ...tree.HookOption) (tree.GroupHook, error) {
	opts = append([]tree.HookOption{
		tree.WithHookName("fail-on-existing-group"),
		tree.WithHookPriority(0),
	}, opts...)

	return &FailOnExistingGroup{tree.NewBaseHook(opts...)}, nil
}
