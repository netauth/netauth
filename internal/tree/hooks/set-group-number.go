package hooks

import (
	"context"

	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

// SetGroupNumber assigns a group number either by using the
// statically provided one, or dynamically requesting one from the
// database.
type SetGroupNumber struct {
	tree.BaseHook
}

// Run will set the group number on g.  If dg.Number is provided as a
// non-zero positive integer, it will be used directly (care should be
// taken this number is not already allocated).  If dg.Number is -1, a
// number will be dynamically provisioned by the database.  It is
// recommended to use automatic provisioning unless strictly necessary
// to do otherwise.
func (s *SetGroupNumber) Run(ctx context.Context, g, dg *pb.Group) error {
	if dg.GetNumber() == -1 {
		number, err := s.Storage().NextGroupNumber(ctx)
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
func NewSetGroupNumber(opts ...tree.HookOption) (tree.GroupHook, error) {
	opts = append([]tree.HookOption{
		tree.WithHookName("set-group-number"),
		tree.WithHookPriority(50),
	}, opts...)

	return &SetGroupNumber{tree.NewBaseHook(opts...)}, nil
}
