package hooks

import (
	"context"

	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

// SetGroupDisplayName copies the display name from dg to g.
type SetGroupDisplayName struct {
	tree.BaseHook
}

// Run copies the DisplayName from dg to g, no checking is performed.
func (*SetGroupDisplayName) Run(_ context.Context, g, dg *pb.Group) error {
	g.DisplayName = dg.DisplayName
	return nil
}

func init() {
	startup.RegisterCallback(setGroupDisplayNameCB)
}

func setGroupDisplayNameCB() {
	tree.RegisterGroupHookConstructor("set-group-displayname", NewSetGroupDisplayName)
}

// NewSetGroupDisplayName returns an initialized hook ready for use.
func NewSetGroupDisplayName(c tree.RefContext) (tree.GroupHook, error) {
	return &SetGroupDisplayName{tree.NewBaseHook("set-group-displayname", 50)}, nil
}
