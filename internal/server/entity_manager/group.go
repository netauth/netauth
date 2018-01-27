package entity_manager

import (
	"log"

	"github.com/NetAuth/NetAuth/pkg/errors"
	pb "github.com/NetAuth/NetAuth/pkg/proto"
)

// newGroup adds a group to the datastore if it does not currently
// exist.  If the group exists then it cannot be added and an error is
// returned.
func (emds *EMDataStore) newGroup(name, displayName string, gidNumber int32) error {
	if _, err := emds.db.LoadGroup(name); err == nil {
		log.Printf("Group '%s' already exists!", name)
		return errors.E_DUPLICATE_GROUP_ID
	}

	if _, err := emds.db.LoadGroupNumber(gidNumber); err == nil || gidNumber == 0 {
		log.Printf("Group number %d is already assigned!", gidNumber)
		return errors.E_DUPLICATE_GROUP_NUMBER
	}

	if gidNumber == -1 {
		var err error
		gidNumber, err = emds.nextGIDNumber()
		if err != nil {
			return err
		}
	}

	newGroup := &pb.Group{
		Name:        &name,
		DisplayName: &displayName,
		GidNumber:   &gidNumber,
	}

	// Save the group
	if err := emds.db.SaveGroup(newGroup); err != nil {
		return err
	}

	log.Printf("Allocated new group '%s'", name)
	return nil
}

func (emds *EMDataStore) NewGroup(requestID, requestSecret, name, displayName string, gidNumber int32) error {
	// Validate that the entity is real and is permitted to
	// perform this action.
	if err := emds.validateEntityCapabilityAndSecret(requestID, requestSecret, "CREATE_GROUP"); err != nil {
		return err
	}

	// Attempt to create the group as specified.
	if err := emds.newGroup(name, displayName, gidNumber); err != nil {
		return err
	}

	return nil
}

// Convenience function to get the nextGIDNumber.  This is very
// inefficient but it only is called when a new group is being
// created, which is hopefully infrequent.
func (emds *EMDataStore) nextGIDNumber() (int32, error) {
	var largest int32 = 0

	l, err := emds.db.DiscoverGroupNames()
	if err != nil {
		return 0, err
	}
	for _, i := range l {
		g, err := emds.db.LoadGroup(i)
		if err != nil {
			return 0, err
		}
		if g.GetGidNumber() > largest {
			largest = g.GetGidNumber()
		}
	}

	return largest + 1, nil
}

// listMembers takes a group ID in and returns a slice of entities
// that are in that group.
func (emds *EMDataStore) listMembers(groupID string) ([]*pb.Entity, error) {
	// 'ALL' is a special groupID which returns everything, this
	// isn't a group that exists in a real sense, it just serves
	// to return a global list as a convenience.
	if groupID == "ALL" {
		var entities []*pb.Entity
		el, err := emds.db.DiscoverEntityIDs()
		if err != nil {
			return nil, err
		}
		for _, en := range el {
			e, err := emds.db.LoadEntity(en)
			if err != nil {
				return nil, err
			}
			entities = append(entities, e)
		}
		return entities, nil
	}

	// No group matched (likely because no other group mechanisms
	// are implemented).
	return nil, errors.E_NO_GROUP
}

// ListMembers fulfills the same function as the private version of
// this function, but with one crucial difference, it produces copies
// of the entities that have the secret redacted.
func (emds *EMDataStore) ListMembers(groupID string) ([]*pb.Entity, error) {
	// This set of members has secrets and can't be returned.
	members, err := emds.listMembers(groupID)
	if err != nil {
		return nil, err
	}

	var safeMembers []*pb.Entity
	for _, e := range members {
		ne, err := safeCopyEntity(e)
		if err != nil {
			return nil, err
		}
		safeMembers = append(safeMembers, ne)
	}

	return safeMembers, nil
}
