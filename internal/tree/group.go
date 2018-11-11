package tree

import (
	"log"
	"strings"

	"github.com/NetAuth/NetAuth/internal/tree/util"
	"github.com/golang/protobuf/proto"

	pb "github.com/NetAuth/Protocol"
)

// NewGroup adds a group to the datastore if it does not currently
// exist.  If the group exists then it cannot be added and an error is
// returned.
func (m *Manager) NewGroup(name, displayName, managedBy string, number int32) error {
	gp := GroupProcessor{
		Group: &pb.Group{},
		RequestData: &pb.Group{
			Name:        &name,
			DisplayName: &displayName,
			ManagedBy:   &managedBy,
			Number:      &number,
		},
	}

	if err := gp.FetchHooks("CREATE-GROUP", m.groupProcesses); err != nil {
		log.Fatal(err)
	}
	_, err := gp.Run()
	return err
}

// GetGroupByName fetches a group by name and returns a pointer to the
// group and a nil error.  If the group cannot be loaded the error
// will explain why.  This is very thin since it just obtains a value
// from the storage layer.
func (m *Manager) GetGroupByName(name string) (*pb.Group, error) {
	return m.db.LoadGroup(name)
}

// DeleteGroup unsurprisingly deletes a group.  There's no real logic
// here, it just passes the delete call through to the storage layer.
func (m *Manager) DeleteGroup(name string) error {
	return m.db.DeleteGroup(name)
}

// UpdateGroupMeta updates metadata within the group.  Certain
// information is not mutable and so that information is not merged
// in.
func (m *Manager) UpdateGroupMeta(name string, update *pb.Group) error {
	g, err := m.GetGroupByName(name)
	if err != nil {
		return err
	}

	// Stash and clear some choice values
	gName := update.GetName()
	number := update.GetNumber()

	update.Name = nil
	update.Number = nil

	proto.Merge(g, update)

	if err := m.db.SaveGroup(g); err != nil {
		return err
	}

	// Put the values back, since this was accessed by pointer
	update.Name = &gName
	update.Number = &number

	return nil
}

// ManageUntypedGroupMeta handles the things that may be annotated
// onto a group.  These annotations should be used sparingly as they
// incur a non-trivial lookup cost on the server.
func (m *Manager) ManageUntypedGroupMeta(name, mode, key, value string) ([]string, error) {
	// Load Entity
	g, err := m.GetGroupByName(name)
	if err != nil {
		return nil, err
	}

	// Patch the KV slice
	tmp := util.PatchKeyValueSlice(g.GetUntypedMeta(), mode, key, value)

	// If this was a read, bail out now with whatever was read
	if strings.ToUpper(mode) == "READ" {
		return tmp, nil
	}

	// Save changes
	g.UntypedMeta = tmp
	if err := m.db.SaveGroup(g); err != nil {
		return nil, err
	}
	return nil, nil
}

// ListGroups literally returns a list of groups
func (m *Manager) ListGroups() ([]*pb.Group, error) {
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

// SetCapability sets a capability on an group.  The set operation is
// idempotent.
func (m *Manager) setGroupCapability(g *pb.Group, c string) error {
	// If no capability was supplied, bail out.
	if len(c) == 0 {
		return ErrUnknownCapability
	}

	cap := pb.Capability(pb.Capability_value[c])

	for _, a := range g.Capabilities {
		if a == cap {
			// The group already has this capability
			// directly, don't add it again.
			return nil
		}
	}

	g.Capabilities = append(g.Capabilities, cap)

	if err := m.db.SaveGroup(g); err != nil {
		return err
	}

	log.Printf("Set capability %s on group '%s'", c, g.GetName())
	return nil
}

// removeCapability removes a capability on an group
func (m *Manager) removeGroupCapability(g *pb.Group, c string) error {
	// If no capability was supplied, bail out.
	if len(c) == 0 {
		return ErrUnknownCapability
	}

	cap := pb.Capability(pb.Capability_value[c])
	var ncaps []pb.Capability

	for _, a := range g.Capabilities {
		if a == cap {
			continue
		}
		ncaps = append(ncaps, a)
	}

	g.Capabilities = ncaps

	if err := m.db.SaveGroup(g); err != nil {
		return err
	}

	log.Printf("Removed capability %s on group '%s'", c, g.GetName())
	return nil
}

// SetGroupCapabilityByName is a convenience function to get the group
// and hand it off to the actual setGroupCapability function
func (m *Manager) SetGroupCapabilityByName(name string, c string) error {
	g, err := m.db.LoadGroup(name)
	if err != nil {
		return err
	}

	return m.setGroupCapability(g, c)
}

// RemoveGroupCapabilityByName is a convenience function to get the group
// and hand it off to the actual removeGroupCapability function
func (m *Manager) RemoveGroupCapabilityByName(name string, c string) error {
	g, err := m.db.LoadGroup(name)
	if err != nil {
		return err
	}

	return m.removeGroupCapability(g, c)
}
