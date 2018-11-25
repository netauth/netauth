package tree

import (
	"fmt"
	"strings"

	"github.com/NetAuth/NetAuth/internal/tree/util"

	pb "github.com/NetAuth/Protocol"
)

// NewGroup adds a group to the datastore if it does not currently
// exist.  If the group exists then it cannot be added and an error is
// returned.
func (m *Manager) NewGroup(name, displayName, managedBy string, number int32) error {
	rg := &pb.Group{
		Name:        &name,
		DisplayName: &displayName,
		ManagedBy:   &managedBy,
		Number:      &number,
	}

	_, err := m.RunGroupChain("CREATE", rg)
	return err
}

// GetGroupByName fetches a group by name and returns a pointer to the
// group and a nil error.  If the group cannot be loaded the error
// will explain why.  This is very thin since it just obtains a value
// from the storage layer.
func (m *Manager) GetGroupByName(name string) (*pb.Group, error) {
	rg := &pb.Group{
		Name: &name,
	}

	return m.RunGroupChain("FETCH", rg)
}

// DeleteGroup unsurprisingly deletes a group.  There's no real logic
// here, it just passes the delete call through to the storage layer.
func (m *Manager) DeleteGroup(name string) error {
	rg := &pb.Group{
		Name: &name,
	}

	_, err := m.RunGroupChain("DESTROY", rg)
	return err
}

// UpdateGroupMeta updates metadata within the group.  Certain
// information is not mutable and so that information is not merged
// in.
func (m *Manager) UpdateGroupMeta(name string, update *pb.Group) error {
	update.Name = &name
	_, err := m.RunGroupChain("MERGE-METADATA", update)
	return err
}

// ManageUntypedGroupMeta handles the things that may be annotated
// onto a group.  These annotations should be used sparingly as they
// incur a non-trivial lookup cost on the server.
func (m *Manager) ManageUntypedGroupMeta(name, mode, key, value string) ([]string, error) {
	rg := &pb.Group{
		Name:        &name,
		UntypedMeta: []string{fmt.Sprintf("%s:%s", key, value)},
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

	g, err := m.RunGroupChain(chain, rg)
	if err != nil {
		return nil, err
	}

	// If this was a read, bail out now with whatever was read
	if strings.ToUpper(mode) == "READ" {
		return util.PatchKeyValueSlice(g.UntypedMeta, "READ", key, ""), nil
	}
	return nil, nil
}

// SetGroupCapabilityByName is a convenience function to get the group
// and hand it off to the actual setGroupCapability function
func (m *Manager) SetGroupCapabilityByName(name string, c string) error {
	capIndex, ok := pb.Capability_value[c]
	if !ok {
		return ErrUnknownCapability
	}

	rg := &pb.Group{
		Name:         &name,
		Capabilities: []pb.Capability{pb.Capability(capIndex)},
	}

	_, err := m.RunGroupChain("SET-CAPABILITY", rg)
	return err
}

// RemoveGroupCapabilityByName is a convenience function to get the group
// and hand it off to the actual removeGroupCapability function
func (m *Manager) RemoveGroupCapabilityByName(name string, c string) error {
	capIndex, ok := pb.Capability_value[c]
	if !ok {
		return ErrUnknownCapability
	}

	rg := &pb.Group{
		Name:         &name,
		Capabilities: []pb.Capability{pb.Capability(capIndex)},
	}

	_, err := m.RunGroupChain("DROP-CAPABILITY", rg)
	return err
}
