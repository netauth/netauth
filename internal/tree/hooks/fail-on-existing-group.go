package hooks

import (
	"github.com/netauth/netauth/internal/db"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

// FailOnExistingGroup is a hook that can be used to guard creation
// processes on groups.
type FailOnExistingGroup struct {
	tree.BaseHook
	db.DB
}

// Run contacts the datastore and attempts to load the group specified
// by dg.  If the group loads successfully then an error is returned,
// in other cases nil is returned.
func (f *FailOnExistingGroup) Run(g, dg *pb.Group) error {
	if _, err := f.LoadGroup(dg.GetName()); err == nil {
		return tree.ErrDuplicateGroupName
	}
	return nil
}

func init() {
	tree.RegisterGroupHookConstructor("fail-on-existing-group", NewFailOnExistingGroup)
}

// NewFailOnExistingGroup returns an initialized hook ready for use.
func NewFailOnExistingGroup(c tree.RefContext) (tree.GroupHook, error) {
	return &FailOnExistingGroup{tree.NewBaseHook("fail-on-existing-group", 0), c.DB}, nil
}
