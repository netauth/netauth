package entity_manager

import (
	pb "github.com/NetAuth/NetAuth/pkg/proto"
)

// addEntityToGroup adds an entity to a group by name, if the entity
// was already in the group the function will return with a nil error.
func (emds *EMDataStore) addEntityToGroup(e *pb.Entity, groupName string) error {
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

// AddEntityToGroup is an external facing function which handles the
// authorization of the change and performs the change.
func (emds *EMDataStore) AddEntityToGroup(mer *pb.ModGroupDirectMembershipRequest) error {
	requestID := mer.GetEntity().GetID()
	requestSecret := mer.GetEntity().GetSecret()

	// Validate that the entity is real and is permitted to
	// perform this action.
	if err := emds.validateEntityCapabilityAndSecret(requestID, requestSecret, "MODIFY_GROUP_MEMBERS"); err != nil {
		return err
	}

	entity, err := emds.db.LoadEntity(mer.GetModEntity())
	if err != nil {
		return err
	}

	if err := emds.addEntityToGroup(entity, mer.GetGroupName()); err != nil {
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

// removeEntityFromGroup removes an entity from the named group.  If
// the entity was not in the group to begin with then nil will be
// returned as the error.
func (emds *EMDataStore) removeEntityFromGroup(e *pb.Entity, groupName string) {
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

// RemoveEntityFromGroup is an external facing function which handles
// the authorization of the change and removes an entity from a group
// in which they have direct membership.
func (emds *EMDataStore) RemoveEntityFromGroup(mer *pb.ModGroupDirectMembershipRequest) error {
	requestID := mer.GetEntity().GetID()
	requestSecret := mer.GetEntity().GetSecret()

	// Validate that the entity is real and is permitted to
	// perform this action.
	if err := emds.validateEntityCapabilityAndSecret(requestID, requestSecret, "MODIFY_GROUP_MEMBERS"); err != nil {
		return err
	}

	entity, err := emds.db.LoadEntity(mer.GetModEntity())
	if err != nil {
		return err
	}

	emds.removeEntityFromGroup(entity, mer.GetGroupName())
	return nil
}
