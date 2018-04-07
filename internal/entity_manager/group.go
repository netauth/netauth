package entity_manager

import (
	"log"

	"github.com/golang/protobuf/proto"

	"github.com/NetAuth/NetAuth/pkg/errors"
	pb "github.com/NetAuth/NetAuth/pkg/proto"
)

// NewGroup adds a group to the datastore if it does not currently
// exist.  If the group exists then it cannot be added and an error is
// returned.
func (emds *EMDataStore) NewGroup(name, displayName string, gidNumber int32) error {
	if _, err := emds.GetGroupByName(name); err == nil {
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

// getGroupByName fetches a group by name and returns a pointer to the
// group and a nil error.  If the group cannot be loaded the error
// will explain why.  This is very thin since it just obtains a value
// from the storage layer.
func (emds *EMDataStore) GetGroupByName(name string) (*pb.Group, error) {
	return emds.db.LoadGroup(name)
}

// deleteGroup unsurprisingly deletes a group.  There's no real logic
// here, it just passes the delete call through to the storage layer.
func (emds *EMDataStore) DeleteGroup(name string) error {
	return emds.db.DeleteGroup(name)
}

// updateGroupMeta updates metadata within the group.  Certain
// information is not mutable and so that information is not merged
// in.
func (emds *EMDataStore) UpdateGroupMeta(name string, update *pb.Group) error {
	g, err := emds.GetGroupByName(name)
	if err != nil {
		return err
	}

	update.Name = nil
	update.GidNumber = nil

	proto.Merge(g, update)

	if err := emds.db.SaveGroup(g); err != nil {
		return err
	}

	return nil
}

// ListGroups literally returns a list of groups
func (emds *EMDataStore) ListGroups() ([]*pb.Group, error) {
	names, err := emds.db.DiscoverGroupNames()
	if err != nil {
		return nil, err
	}

	groups := []*pb.Group{}
	for _, name := range names {
		g, err := emds.db.LoadGroup(name)
		if err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil
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
	var entities []*pb.Entity

	// 'ALL' is a special groupID which returns everything, this
	// isn't a group that exists in a real sense, it just serves
	// to return a global list as a convenience.
	if groupID == "ALL" {
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

	// If its not the all group then we check to make sure the
	// group exists at all
	if _, err := emds.db.LoadGroup(groupID); err != nil {
		return nil, err
	}

	// Now we can be reasonably sure the group exists, this next
	// bit is stupidly inefficient, but is the only way to extract
	// the members since the membership graph has the arrows going
	// the other way.
	el, err := emds.listMembers("ALL")
	if err != nil {
		return nil, err
	}
	for _, e := range el {
		for _, g := range emds.GetDirectGroups(e) {
			if g == groupID {
				entities = append(entities, e)
			}
		}
	}
	return entities, nil
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
