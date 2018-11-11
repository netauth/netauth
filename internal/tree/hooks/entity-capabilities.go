package hooks

import (
	"github.com/NetAuth/NetAuth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

type ManageEntityCapabilities struct {
	tree.BaseHook
	mode bool
}

func (mec *ManageEntityCapabilities) Run(e, de *pb.Entity) error {
	for _, cap := range de.Meta.Capabilities {
		if mec.mode {
			// Add mode
			de.Meta.Capabilities = mec.addCapability(cap, de.Meta.Capabilities)
		} else {
			// Del mode
			de.Meta.Capabilities = mec.delCapability(cap, de.Meta.Capabilities)
		}
	}
	return nil
}

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
	ncaps = []pb.Capability{cap}

	return ncaps
}

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

func NewSetEntityCapability(c tree.RefContext) (tree.EntityProcessorHook, error) {
	return &ManageEntityCapabilities{tree.NewBaseHook("set-entity-capability", 50), true}, nil
}

func NewRemoveEntityCapability(c tree.RefContext) (tree.EntityProcessorHook, error) {
	return &ManageEntityCapabilities{tree.NewBaseHook("remove-entity-capability", 50), false}, nil
}
