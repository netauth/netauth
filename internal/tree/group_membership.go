package tree

import (
	pb "github.com/NetAuth/NetAuth/pkg/proto"
)

// AddEntityToGroup is the same as the internal function, but takes an
// entity ID rather than a pointer
func (m Manager) AddEntityToGroup(entityID, groupName string) error {
	e, err := m.db.LoadEntity(entityID)
	if err != nil {
		return err
	}
	return m.addEntityToGroup(e, groupName)
}

// addEntityToGroup adds an entity to a group by name, if the entity
// was already in the group the function will return with a nil error.
func (m Manager) addEntityToGroup(e *pb.Entity, groupName string) error {
	if _, err := m.db.LoadGroup(groupName); err != nil {
		return err
	}

	if e.GetMeta() == nil {
		e.Meta = &pb.EntityMeta{}
	}

	// First we check if the entity is a member of the group
	// directly.
	groupNames := e.GetMeta().GetGroups()
	for _, g := range groupNames {
		if g == groupName {
			return nil
		}
	}

	// At this point we can be reasonably certain that the entity
	// is not in the named group via direct membership.
	e.Meta.Groups = append(e.Meta.Groups, groupName)

	if err := m.db.SaveEntity(e); err != nil {
		return err
	}
	return nil
}

// GetDirectGroups gets the direct groups of an entity.
func (m Manager) GetDirectGroups(e *pb.Entity) []string {
	if e.GetMeta() == nil {
		return []string{}
	}

	return e.GetMeta().GetGroups()
}

// RemoveEntityFromGroup performs the same function as the internal
// variant, but does so by name rather than by entity pointer.
func (m Manager) RemoveEntityFromGroup(entityID, groupName string) error {
	e, err := m.db.LoadEntity(entityID)
	if err != nil {
		return err
	}
	m.removeEntityFromGroup(e, groupName)
	return nil
}

// removeEntityFromGroup removes an entity from the named group.  If
// the entity was not in the group to begin with then nil will be
// returned as the error.
func (m Manager) removeEntityFromGroup(e *pb.Entity, groupName string) {
	if e.GetMeta() == nil {
		return
	}

	newGroups := []string{}
	for _, g := range e.GetMeta().GetGroups() {
		if g == groupName {
			continue
		}
		newGroups = append(newGroups, g)
	}
	e.Meta.Groups = newGroups

	if err := m.db.SaveEntity(e); err != nil {
		return
	}
	return
}
