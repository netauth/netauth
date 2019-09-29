package rpc2

import (
	types "github.com/NetAuth/Protocol"
)

func (s *Server) getCapabilitiesForEntity(id string) []types.Capability {
	// Get the full fledged entity; we can assert no error here
	// since the entity was just loaded to perform an
	// authentication check prior to calling this function.
	// There's a minimal risk that the entity has vanished, but if
	// that's the case then this will return empty since the
	// FetchEntity function will return an empty entity, thus no
	// authorization attack can be performed via this vector.
	e, _ := s.FetchEntity(id)

	// First get the capabilities that are provided by the entity
	// itself.
	caps := make(map[types.Capability]int)
	if e.GetMeta() != nil {
		for _, c := range e.GetMeta().GetCapabilities() {
			caps[c]++
		}
	}

	// Next get the capabilities that are provided by any groups
	// the entity may be in; include indirects for authentication
	// queries.  We can assert no error here because the worst
	// case is a group is returned that is completely empty.  In
	// this case there will simply be no additional capabilities
	// granted, and this function can continue without incident.
	groupNames := s.GetMemberships(e, true)
	for _, name := range groupNames {
		g, _ := s.FetchGroup(name)
		for _, c := range g.GetCapabilities() {
			caps[c]++
		}
	}

	// Flatten the capabilities out into a list
	capabilities := make([]types.Capability, len(caps))
	i := 0
	for c := range caps {
		capabilities[i] = c
		i++
	}
	return capabilities
}
