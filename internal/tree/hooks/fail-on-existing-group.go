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
	tree.DB
}

// Run contacts the datastore and attempts to load the group specified
// by dg.  If the group loads successfully then an error is returned,
// in other cases nil is returned.
func (f *FailOnExistingGroup) Run(ctx context.Context, g, dg *pb.Group) error {
	if _, err := f.LoadGroup(ctx, dg.GetName()); err == nil {
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
func NewFailOnExistingGroup(c tree.RefContext) (tree.GroupHook, error) {
	return &FailOnExistingGroup{tree.NewBaseHook("fail-on-existing-group", 0), c.DB}, nil
}
