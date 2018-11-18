package tree

import (
	"fmt"
	"log"
	"strings"

	"github.com/NetAuth/NetAuth/internal/tree/util"

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

	if err := gp.FetchHooks("CREATE", m.groupProcesses); err != nil {
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
	gp := GroupProcessor{
		Group: &pb.Group{},
		RequestData: &pb.Group{
			Name: &name,
		},
	}

	if err := gp.FetchHooks("FETCH", m.groupProcesses); err != nil {
		log.Fatal(err)
	}
	return gp.Run()
}

// DeleteGroup unsurprisingly deletes a group.  There's no real logic
// here, it just passes the delete call through to the storage layer.
func (m *Manager) DeleteGroup(name string) error {
	gp := GroupProcessor{
		Group: &pb.Group{},
		RequestData: &pb.Group{
			Name: &name,
		},
	}

	if err := gp.FetchHooks("DESTROY", m.groupProcesses); err != nil {
		log.Fatal(err)
	}
	_, err := gp.Run()
	return err
}

// UpdateGroupMeta updates metadata within the group.  Certain
// information is not mutable and so that information is not merged
// in.
func (m *Manager) UpdateGroupMeta(name string, update *pb.Group) error {
	gp := GroupProcessor{
		Group:       &pb.Group{},
		RequestData: update,
	}

	if err := gp.FetchHooks("MERGE-METADATA", m.groupProcesses); err != nil {
		log.Fatal(err)
	}
	_, err := gp.Run()
	return err
}

// ManageUntypedGroupMeta handles the things that may be annotated
// onto a group.  These annotations should be used sparingly as they
// incur a non-trivial lookup cost on the server.
func (m *Manager) ManageUntypedGroupMeta(name, mode, key, value string) ([]string, error) {
	gp := GroupProcessor{
		Group: &pb.Group{},
		RequestData: &pb.Group{
			Name:        &name,
			UntypedMeta: []string{fmt.Sprintf("%s:%s", key, value)},
		},
	}

	// Mode switch and select appropriate processor chain.
	chain := "FETCH"
	switch strings.ToUpper(mode) {
	case "UPSERT":
		chain = "UGM-UPSERT"
	case "CLEARFUZZY":
		chain = "UGM-CLEARFUZZY"
	case "CLEAREXACT":
		chain = "UGM-CLEAREXACT"
	default:
		mode = "READ"
	}

	// Process transaction
	if err := gp.FetchHooks(chain, m.groupProcesses); err != nil {
		log.Fatal(err)
	}
	g, err := gp.Run()
	if err != nil {
		return nil, err
	}

	// If this was a read, bail out now with whatever was read
	if strings.ToUpper(mode) == "READ" {
		return util.PatchKeyValueSlice(g.UntypedMeta, "READ", key, ""), nil
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

// SetGroupCapabilityByName is a convenience function to get the group
// and hand it off to the actual setGroupCapability function
func (m *Manager) SetGroupCapabilityByName(name string, c string) error {
	capIndex, ok := pb.Capability_value[c]
	if !ok {
		return ErrUnknownCapability
	}

	gp := GroupProcessor{
		Group: &pb.Group{},
		RequestData: &pb.Group{
			Name:         &name,
			Capabilities: []pb.Capability{pb.Capability(capIndex)},
		},
	}

	if err := gp.FetchHooks("SET-CAPABILITY", m.groupProcesses); err != nil {
		log.Fatal(err)
	}
	_, err := gp.Run()
	return err
}

// RemoveGroupCapabilityByName is a convenience function to get the group
// and hand it off to the actual removeGroupCapability function
func (m *Manager) RemoveGroupCapabilityByName(name string, c string) error {
	capIndex, ok := pb.Capability_value[c]
	if !ok {
		return ErrUnknownCapability
	}

	gp := GroupProcessor{
		Group: &pb.Group{},
		RequestData: &pb.Group{
			Name:         &name,
			Capabilities: []pb.Capability{pb.Capability(capIndex)},
		},
	}

	if err := gp.FetchHooks("DROP-CAPABILITY", m.groupProcesses); err != nil {
		log.Fatal(err)
	}
	_, err := gp.Run()
	return err
}
