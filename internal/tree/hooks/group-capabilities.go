package hooks

import (
	"context"

	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

// ManageGroupCapabilities is a configurable runtime hook that adds
// or removes capabilities as configured.  Two convenience
// constructors exist to return hooks in either mode.
type ManageGroupCapabilities struct {
	tree.BaseHook
	mode bool
}

// Run modifies the stored capabilities on an entity depending on the
// value of the mode variable.  When the mode is set to true, any
// capabilities stored in de will be copied to e if they are not
// already present.  In false capabilities will be subtracted.
func (mec *ManageGroupCapabilities) Run(_ context.Context, g, dg *pb.Group) error {
	if len(dg.Capabilities) == 0 {
		return tree.ErrUnknownCapability
	}
	for _, cap := range dg.Capabilities {
		if mec.mode {
			// Add mode
			g.Capabilities = addCapability(cap, g.Capabilities)
		} else {
			// Del mode
			g.Capabilities = delCapability(cap, g.Capabilities)
		}
	}
	return nil
}

func init() {
	startup.RegisterCallback(groupCapabilitiesCB)
}

func groupCapabilitiesCB() {
	tree.RegisterGroupHookConstructor("set-group-capability", NewSetGroupCapability)
	tree.RegisterGroupHookConstructor("remove-group-capability", NewRemoveGroupCapability)
}

// NewSetGroupCapability returns a ManageGroupCapability hook
// pre-configured into the additive mode.
func NewSetGroupCapability(opts ...tree.HookOption) (tree.GroupHook, error) {
	opts = append([]tree.HookOption{
		tree.WithHookName("set-group-capability"),
		tree.WithHookPriority(50),
	}, opts...)

	return &ManageGroupCapabilities{tree.NewBaseHook(opts...), true}, nil
}

// NewRemoveGroupCapability returns a ManageGroupCapability hook
// pre-configured into the subtractive mode.s
func NewRemoveGroupCapability(opts ...tree.HookOption) (tree.GroupHook, error) {
	opts = append([]tree.HookOption{
		tree.WithHookName("remove-group-capability"),
		tree.WithHookPriority(50),
	}, opts...)

	return &ManageGroupCapabilities{tree.NewBaseHook(opts...), false}, nil
}
