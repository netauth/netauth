package hooks

import (
	"context"

	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

// SetEntityNumber sets the number on an entity either statically with
// the number provided or dynamically with an automatically chosen
// number.
type SetEntityNumber struct {
	tree.BaseHook
}

// Run will provision a number in one of two ways.  If the number is
// not equal to -1 then it will be used directly with no further
// checks and will be applied to the entity.  If the number is -1 then
// the data storage system will be queried for the next available
// number.  These numbers are not guaranteed to be in order or have
// any mathematical progression, only uniqueness.
func (s *SetEntityNumber) Run(ctx context.Context, e, de *pb.Entity) error {
	if de.GetNumber() == -1 {
		n, err := s.Storage().NextEntityNumber(ctx)
		if err != nil {
			return err
		}
		e.Number = &n
		return nil
	}
	e.Number = de.Number
	return nil
}

func init() {
	startup.RegisterCallback(setEntityNumberCB)
}

func setEntityNumberCB() {
	tree.RegisterEntityHookConstructor("set-entity-number", NewSetEntityNumber)
}

// NewSetEntityNumber returns a SetEntityNumber hook ready for use.
func NewSetEntityNumber(opts ...tree.HookOption) (tree.EntityHook, error) {
	opts = append([]tree.HookOption{
		tree.WithHookName("set-entity-number"),
		tree.WithHookPriority(50),
	}, opts...)

	return &SetEntityNumber{tree.NewBaseHook(opts...)}, nil
}
