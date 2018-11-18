package hooks

import (
	"github.com/NetAuth/NetAuth/internal/tree"

	pb "github.com/NetAuth/Protocol"
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
func (mec *ManageGroupCapabilities) Run(g, dg *pb.Group) error {
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
	tree.RegisterGroupHookConstructor("set-group-capability", NewSetGroupCapability)
	tree.RegisterGroupHookConstructor("remove-group-capability", NewRemoveGroupCapability)
}

// NewSetGroupCapability returns a ManageGroupCapability hook
// pre-configured into the additive mode.
func NewSetGroupCapability(c tree.RefContext) (tree.GroupProcessorHook, error) {
	return &ManageGroupCapabilities{tree.NewBaseHook("set-group-capability", 50), true}, nil
}

// NewRemoveGroupCapability returns a ManageGroupCapability hook
// pre-configured into the subtractive mode.s
func NewRemoveGroupCapability(c tree.RefContext) (tree.GroupProcessorHook, error) {
	return &ManageGroupCapabilities{tree.NewBaseHook("remove-group-capability", 50), false}, nil
}
