package entity_manager

import (
	pb "github.com/NetAuth/NetAuth/pkg/proto"
)

// addEntityToDirectGroup adds an entity to a group by name, if the entity
// was already in the group the function will return with a nil error.
func (emds *EMDataStore) addEntityToDirectGroup(e *pb.Entity, groupName string) error {
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

	if err := emds.db.SaveEntity(e); err != nil {
		return err
	}
	return nil
}

// getDirectGroups gets the direct groups of an entity.
func (emds *EMDataStore) getDirectGroups(e *pb.Entity) []string {
	if e.GetMeta() == nil {
		return []string{}
	}

	return e.GetMeta().GetGroups()
}

// removeEntityFromDirectGroup removes  an entity from the  named group.  If
// the entity  was not  in the group  to begin with  then nil  will be
// returned as the error.
func (emds *EMDataStore) removeEntityFromDirectGroup(e *pb.Entity, groupName string) {
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
}
