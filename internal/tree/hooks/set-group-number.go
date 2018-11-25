package hooks

import (
	"github.com/NetAuth/NetAuth/internal/db"
	"github.com/NetAuth/NetAuth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

// SetGroupNumber assigns a group number either by using the
// statically provided one, or dynamically requesting one from the
// database.
type SetGroupNumber struct {
	tree.BaseHook
	db.DB
}

// Run will set the group number on g.  If dg.Number is provided as a
// non-zero positive integer, it will be used directly (care should be
// taken this number is not already allocated).  If dg.Number is -1, a
// number will be dynamically provisioned by the database.  It is
// recommended to use automatic provisioning unless strictly necessary
// to do otherwise.
func (s *SetGroupNumber) Run(g, dg *pb.Group) error {
	if dg.GetNumber() == -1 {
		number, err := s.NextGroupNumber()
		if err != nil {
			return err
		}
		g.Number = &number
		return nil
	}
	g.Number = dg.Number
	return nil
}

func init() {
	tree.RegisterGroupHookConstructor("set-group-number", NewSetGroupNumber)
}

// NewSetGroupNumber returns a hook initialized and ready for use.
func NewSetGroupNumber(c tree.RefContext) (tree.GroupHook, error) {
	return &SetGroupNumber{tree.NewBaseHook("set-group-number", 50), c.DB}, nil
}
