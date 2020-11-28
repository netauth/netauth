package hooks

import (
	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

// SetGroupNumber assigns a group number either by using the
// statically provided one, or dynamically requesting one from the
// database.
type SetGroupNumber struct {
	tree.BaseHook
	tree.DB
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
	startup.RegisterCallback(setGroupNumberCB)
}

func setGroupNumberCB() {
	tree.RegisterGroupHookConstructor("set-group-number", NewSetGroupNumber)
}

// NewSetGroupNumber returns a hook initialized and ready for use.
func NewSetGroupNumber(c tree.RefContext) (tree.GroupHook, error) {
	return &SetGroupNumber{tree.NewBaseHook("set-group-number", 50), c.DB}, nil
}
