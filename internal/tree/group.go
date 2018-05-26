package tree

import (
	"log"

	"github.com/golang/protobuf/proto"

	"github.com/NetAuth/NetAuth/pkg/errors"
	pb "github.com/NetAuth/Protocol"
)

// NewGroup adds a group to the datastore if it does not currently
// exist.  If the group exists then it cannot be added and an error is
// returned.
func (m Manager) NewGroup(name, displayName, managedBy string, gidNumber int32) error {
	if _, err := m.GetGroupByName(name); err == nil {
		log.Printf("Group '%s' already exists!", name)
		return errors.E_DUPLICATE_GROUP_ID
	}

	// Verify that the managing group exists.
	if _, err := m.GetGroupByName(managedBy); managedBy != "" && managedBy != name && err != nil {
		return err
	}

	if _, err := m.db.LoadGroupNumber(gidNumber); err == nil || gidNumber == 0 {
		log.Printf("Group number %d is already assigned!", gidNumber)
		return errors.E_DUPLICATE_GROUP_NUMBER
	}

	if gidNumber == -1 {
		var err error
		gidNumber, err = m.nextGIDNumber()
		if err != nil {
			return err
		}
	}

	newGroup := &pb.Group{
		Name:        &name,
		DisplayName: &displayName,
		GidNumber:   &gidNumber,
		ManagedBy:   &managedBy,
	}

	// Save the group
	if err := m.db.SaveGroup(newGroup); err != nil {
		return err
	}

	log.Printf("Allocated new group '%s'", name)
	return nil
}

// getGroupByName fetches a group by name and returns a pointer to the
// group and a nil error.  If the group cannot be loaded the error
// will explain why.  This is very thin since it just obtains a value
// from the storage layer.
func (m Manager) GetGroupByName(name string) (*pb.Group, error) {
	return m.db.LoadGroup(name)
}

// deleteGroup unsurprisingly deletes a group.  There's no real logic
// here, it just passes the delete call through to the storage layer.
func (m Manager) DeleteGroup(name string) error {
	return m.db.DeleteGroup(name)
}

// updateGroupMeta updates metadata within the group.  Certain
// information is not mutable and so that information is not merged
// in.
func (m Manager) UpdateGroupMeta(name string, update *pb.Group) error {
	g, err := m.GetGroupByName(name)
	if err != nil {
		return err
	}

	// Stash and clear some choice values
	gName := update.GetName()
	number := update.GetGidNumber()

	update.Name = nil
	update.GidNumber = nil

	proto.Merge(g, update)

	if err := m.db.SaveGroup(g); err != nil {
		return err
	}

	// Put the values back, since this was accessed by pointer
	update.Name = &gName
	update.GidNumber = &number

	return nil
}

// ListGroups literally returns a list of groups
func (m Manager) ListGroups() ([]*pb.Group, error) {
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

// Convenience function to get the nextGIDNumber.  This is very
// inefficient but it only is called when a new group is being
// created, which is hopefully infrequent.
func (m Manager) nextGIDNumber() (int32, error) {
	var largest int32 = 0

	l, err := m.db.DiscoverGroupNames()
	if err != nil {
		return 0, err
	}
	for _, i := range l {
		g, err := m.db.LoadGroup(i)
		if err != nil {
			return 0, err
		}
		if g.GetGidNumber() > largest {
			largest = g.GetGidNumber()
		}
	}

	return largest + 1, nil
}
