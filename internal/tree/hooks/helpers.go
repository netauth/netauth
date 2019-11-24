package hooks

import (
	"strings"

	pb "github.com/netauth/protocol"
)

func splitKeyValue(s string) (string, string) {
	parts := strings.SplitN(s, ":", 2)
	return parts[0], parts[1]
}

// addCapability is an internal convenience function to add
// capabilities if they do not already exist in a capability slice.
func addCapability(cap pb.Capability, caps []pb.Capability) []pb.Capability {
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
func delCapability(cap pb.Capability, caps []pb.Capability) []pb.Capability {
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
