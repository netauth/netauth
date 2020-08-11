package hooks

import (
	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
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
	startup.RegisterCallback(setGroupNameCB)
}

func setGroupNameCB() {
	tree.RegisterGroupHookConstructor("set-group-name", NewSetGroupName)
}

// NewSetGroupName returns an initialized hook.
func NewSetGroupName(c tree.RefContext) (tree.GroupHook, error) {
	return &SetGroupName{tree.NewBaseHook("set-group-name", 50)}, nil
}
