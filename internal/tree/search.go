package tree

import (
	pb "github.com/NetAuth/Protocol"
)

// SearchGroups returns a list of groups filtered by the search
// criteria.
func (m *Manager) SearchGroups() ([]*pb.Group, error) {
	names, err := m.db.DiscoverGroupNames()
	if err != nil {
		return nil, err
	}

	groups := []*pb.Group{}
	for _, name := range names {
		g, err := m.db.LoadGroup(name)
		if err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil
}

// SearchEntities returns a list of entities filtered by the search
// criteria.
func (m *Manager) SearchEntities() ([]*pb.Entity, error) {
	var entities []*pb.Entity
	el, err := m.db.DiscoverEntityIDs()
	if err != nil {
		return nil, err
	}
	for _, en := range el {
		e, err := m.db.LoadEntity(en)
		if err != nil {
			return nil, err
		}
		entities = append(entities, e)
	}
	return entities, nil
}
