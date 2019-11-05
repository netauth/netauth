package hooks

import (
	"github.com/netauth/netauth/internal/db"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

// SetManagingGroup performs validation checks on the managing group
// and then sets it.
type SetManagingGroup struct {
	tree.BaseHook
	db.DB
}

// Run will attempt to set the managing group of g to the specified
// group on dg.  If the managing group is the empty string,
// i.e. unmanaged, the hook will return immediately, otherwise the
// group is checked for either existence, or identity to the group
// being created.
func (c *SetManagingGroup) Run(g, dg *pb.Group) error {
	// If the managedby field is blank, this group is unmanaged
	// and requires token authority to alter later.
	if dg.GetManagedBy() == "" {
		return nil
	}

	// If the group that is managing this one is the same name
	// (i.e. self-managed) then we return ok regardless of if the
	// group exists in the data store or not.
	if dg.GetName() == dg.GetManagedBy() {
		g.ManagedBy = dg.ManagedBy
		return nil
	}

	// If the group is not self managed but does have a manage by,
	// then the managedby group must exist already.
	if _, err := c.LoadGroup(dg.GetManagedBy()); err != nil {
		return err
	}

	// All must be okay at this point
	g.ManagedBy = dg.ManagedBy
	return nil
}

func init() {
	tree.RegisterGroupHookConstructor("set-managing-group", NewSetManagingGroup)
}

// NewSetManagingGroup returns a hook initialized for use.
func NewSetManagingGroup(c tree.RefContext) (tree.GroupHook, error) {
	return &SetManagingGroup{tree.NewBaseHook("set-managing-group", 10), c.DB}, nil
}
