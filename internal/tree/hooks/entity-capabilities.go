package hooks

import (
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

// ManageEntityCapabilities is a configurable runtime hook that adds
// or removes capabilities as configured.  Two convenience
// constructors exist to return hooks in either mode.
type ManageEntityCapabilities struct {
	tree.BaseHook
	mode bool
}

// Run modifies the stored capabilities on an entity depending on the
// value of the mode variable.  When the mode is set to true, any
// capabilities stored in de will be copied to e if they are not
// already present.  In false capabilities will be subtracted.
func (mec *ManageEntityCapabilities) Run(e, de *pb.Entity) error {
	if de.Meta == nil || len(de.Meta.Capabilities) == 0 {
		return tree.ErrUnknownCapability
	}
	for _, cap := range de.Meta.Capabilities {
		if mec.mode {
			// Add mode
			e.Meta.Capabilities = addCapability(cap, e.Meta.Capabilities)
		} else {
			// Del mode
			e.Meta.Capabilities = delCapability(cap, e.Meta.Capabilities)
		}
	}
	return nil
}

func init() {
	tree.RegisterEntityHookConstructor("set-entity-capability", NewSetEntityCapability)
	tree.RegisterEntityHookConstructor("remove-entity-capability", NewRemoveEntityCapability)
}

// NewSetEntityCapability returns a ManageEntityCapability hook
// pre-configured into the additive mode.
func NewSetEntityCapability(c tree.RefContext) (tree.EntityHook, error) {
	return &ManageEntityCapabilities{tree.NewBaseHook("set-entity-capability", 50), true}, nil
}

// NewRemoveEntityCapability returns a ManageEntityCapability hook
// pre-configured into the subtractive mode.s
func NewRemoveEntityCapability(c tree.RefContext) (tree.EntityHook, error) {
	return &ManageEntityCapabilities{tree.NewBaseHook("remove-entity-capability", 50), false}, nil
}
