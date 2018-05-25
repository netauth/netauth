package rpc

func (s *NetAuthServer) manageByMembership(entityID, groupName string) bool {
	g, err := s.Tree.GetGroupByName(groupName)
	if err != nil {
		// If the group can't be summoned, pessimistically
		// return false
		return false
	}

	// management by group membership is only available if the
	// group is configured to trust another group for this task,
	// so if this is cleared then no group is trusted.
	if g.GetManagedBy() == "" {
		// This group doesn't have delegated administrative
		// properties.
		return false
	}

	// Get the entity itself for a group check
	e, err := s.Tree.GetEntity(entityID)
	if err != nil {
		return false
	}

	// Always include indirects when evaluating if in an
	// administrative group
	groups := s.Tree.GetMemberships(e, true)

	// Check if any of the groups are the one that grants this
	// power
	for _, name := range groups {
		if name == groupName {
			return true
		}
	}

	// Group checks fall through, return false
	return false
}
