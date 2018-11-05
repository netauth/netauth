package hooks

import (
	pb "github.com/NetAuth/Protocol"
)

type RemoveEntityCapability struct{}

func (*RemoveEntityCapability) Name() string  { return "remove-entity-capability" }
func (*RemoveEntityCapability) Priority() int { return 50 }
func (*RemoveEntityCapability) Run(e, de *pb.Entity) error {
	cap := de.GetMeta().GetCapabilities()[0]
	var ncaps []pb.Capability

	for _, a := range e.Meta.Capabilities {
		if a == cap {
			continue
		}
		ncaps = append(ncaps, a)
	}

	// First capability so the loop above didn't iterate
	if len(e.Meta.Capabilities) == 0 {
		ncaps = []pb.Capability{cap}
	}

	e.Meta.Capabilities = ncaps
	return nil
}
