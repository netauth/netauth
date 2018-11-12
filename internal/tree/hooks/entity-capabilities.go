package hooks

import (
	"github.com/NetAuth/NetAuth/internal/tree"

	pb "github.com/NetAuth/Protocol"
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
	for _, cap := range de.Meta.Capabilities {
		if mec.mode {
			// Add mode
			e.Meta.Capabilities = mec.addCapability(cap, e.Meta.Capabilities)
		} else {
			// Del mode
			e.Meta.Capabilities = mec.delCapability(cap, e.Meta.Capabilities)
		}
	}
	return nil
}

// addCapability is an internal convenience function to add
// capabilities if they do not already exist in a capability slice.
func (mec *ManageEntityCapabilities) addCapability(cap pb.Capability, caps []pb.Capability) []pb.Capability {
	var ncaps []pb.Capability

	// Check to make sure that the capability isn't already set
	for _, a := range caps {
		if a == cap {
			return caps
		}
		ncaps = append(ncaps, a)
	}

	// Add the new capability to the list.
	ncaps = append(ncaps, cap)

	return ncaps
}

// delCapability is an internal convenience function to remove
// capabilities that exist in a slice.
func (mec *ManageEntityCapabilities) delCapability(cap pb.Capability, caps []pb.Capability) []pb.Capability {
	var ncaps []pb.Capability
	for _, a := range caps {
		if a == cap {
			// Don't copy the same capability
			continue
		}
		ncaps = append(ncaps, a)
	}
	return ncaps
}

func init() {
	tree.RegisterEntityHookConstructor("set-entity-capability", NewSetEntityCapability)
	tree.RegisterEntityHookConstructor("remove-entity-capability", NewRemoveEntityCapability)
}

// NewSetEntityCapability returns a ManageEntityCapability hook
// pre-configured into the additive mode.
func NewSetEntityCapability(c tree.RefContext) (tree.EntityProcessorHook, error) {
	return &ManageEntityCapabilities{tree.NewBaseHook("set-entity-capability", 50), true}, nil
}

// NewRemoveEntityCapability returns a ManageEntityCapability hook
// pre-configured into the subtractive mode.s
func NewRemoveEntityCapability(c tree.RefContext) (tree.EntityProcessorHook, error) {
	return &ManageEntityCapabilities{tree.NewBaseHook("remove-entity-capability", 50), false}, nil
}
