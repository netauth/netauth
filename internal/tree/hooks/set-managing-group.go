package hooks

import (
	"github.com/NetAuth/NetAuth/internal/db"

	pb "github.com/NetAuth/Protocol"
)

type SetManagingGroup struct {
	db.DB
}

func (*SetManagingGroup) Name() string  { return "check-managing-group" }
func (*SetManagingGroup) Priority() int { return 10 }
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
