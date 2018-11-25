package hooks

import (
	"github.com/NetAuth/NetAuth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

// SetGroupDisplayName copies the display name from dg to g.
type SetGroupDisplayName struct {
	tree.BaseHook
}

// Run copies the DisplayName from dg to g, no checking is performed.
func (*SetGroupDisplayName) Run(g, dg *pb.Group) error {
	g.DisplayName = dg.DisplayName
	return nil
}

func init() {
	tree.RegisterGroupHookConstructor("set-group-displayname", NewSetGroupDisplayName)
}

// NewSetGroupDisplayName returns an initialized hook ready for use.
func NewSetGroupDisplayName(c tree.RefContext) (tree.GroupHook, error) {
	return &SetGroupDisplayName{tree.NewBaseHook("set-group-displayname", 50)}, nil
}
