package hooks

import (
	"github.com/NetAuth/NetAuth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

// SetGroupName copies the name from dg to g.
type SetGroupName struct {
	tree.BaseHook
}

// Run sets the name on g to the name on dg, no checks or validation
// are run during this hook.
func (*SetGroupName) Run(g, dg *pb.Group) error {
	g.Name = dg.Name
	return nil
}

func init() {
	tree.RegisterGroupHookConstructor("set-group-name", NewSetGroupName)
}

// NewSetGroupName returns an initialized hook.
func NewSetGroupName(c tree.RefContext) (tree.GroupProcessorHook, error) {
	return &SetGroupName{tree.NewBaseHook("set-group-name", 50)}, nil
}
