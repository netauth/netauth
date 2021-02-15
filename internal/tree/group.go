package tree

import (
	"fmt"
	"strings"

	"github.com/netauth/netauth/internal/tree/util"

	pb "github.com/netauth/protocol"
)

// CreateGroup adds a group to the datastore if it does not currently
// exist.  If the group exists then it cannot be added and an error is
// returned.
func (m *Manager) CreateGroup(name, displayName, managedBy string, number int32) error {
	rg := &pb.Group{
		Name:        &name,
		DisplayName: &displayName,
		ManagedBy:   &managedBy,
		Number:      &number,
	}

	_, err := m.RunGroupChain("CREATE", rg)
	return err
}

// FetchGroup fetches a group by name and returns a pointer to the
// group and a nil error.  If the group cannot be loaded the error
// will explain why.  This is very thin since it just obtains a value
// from the storage layer.
func (m *Manager) FetchGroup(name string) (*pb.Group, error) {
	rg := &pb.Group{
		Name: &name,
	}

	return m.RunGroupChain("FETCH", rg)
}

// DestroyGroup unsurprisingly deletes a group.  There's no real logic
// here, it just passes the delete call through to the storage layer.
func (m *Manager) DestroyGroup(name string) error {
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

// GroupKVGet returns an existing key from a group.  If the key does
// not exist an error is returned.
func (m *Manager) GroupKVGet(name string, keys []*pb.KVData) ([]*pb.KVData, error) {
	g, err := m.FetchGroup(name)
	if err != nil {
		return nil, err
	}

	out := []*pb.KVData{}
	for _, haystack := range g.GetKV() {
		for _, needle := range keys {
			if haystack.GetKey() != needle.GetKey() {
				continue
			}
			out = append(out, haystack)
		}
	}
	if len(out) == 0 {
		return nil, ErrNoSuchKey
	}
	return out, nil
}

// GroupKVAdd adds a new key to a group.  If the key already exists
// an error is returned.
func (m *Manager) GroupKVAdd(name string, d []*pb.KVData) error {
	dg := &pb.Group{
		Name: &name,
		KV:   d,
	}

	_, err := m.RunGroupChain("KV-ADD", dg)
	return err
}

// GroupKVDel removes an existing key from a group.  If the key does
// not exist an error is returned.
func (m *Manager) GroupKVDel(name string, d []*pb.KVData) error {
	dg := &pb.Group{
		Name: &name,
		KV:   d,
	}

	_, err := m.RunGroupChain("KV-DEL", dg)
	return err
}

// GroupKVReplace replaces an existing key on a group.  If the key
// does not exist an error is returned.
func (m *Manager) GroupKVReplace(name string, d []*pb.KVData) error {
	dg := &pb.Group{
		Name: &name,
		KV:   d,
	}

	_, err := m.RunGroupChain("KV-REPLACE", dg)
	return err
}

// SetGroupCapability adds a capability to an existing group.  It
// should be preferred to add capabilities to groups rather than to
// entities directly.
func (m *Manager) SetGroupCapability(name string, c string) error {
	capIndex, ok := pb.Capability_value[c]
	if !ok {
		return ErrUnknownCapability
	}

	cap := pb.Capability(capIndex)
	return m.SetGroupCapability2(name, &cap)
}

// SetGroupCapability2 adds a capability to an existing group, and
// does so with a strongly typed capability pointer.  It should be
// preferred to add capabilities to groups rather than to entities
// directly.
func (m *Manager) SetGroupCapability2(name string, c *pb.Capability) error {
	if c == nil {
		return ErrUnknownCapability
	}

	rg := &pb.Group{
		Name:         &name,
		Capabilities: []pb.Capability{*c},
	}

	_, err := m.RunGroupChain("SET-CAPABILITY", rg)
	return err
}

// DropGroupCapability drops a capability from an existing group.
func (m *Manager) DropGroupCapability(name string, c string) error {
	capIndex, ok := pb.Capability_value[c]
	if !ok {
		return ErrUnknownCapability
	}

	cap := pb.Capability(capIndex)
	return m.DropGroupCapability2(name, &cap)
}

// DropGroupCapability2 drops a capability from an existing group, and
// does so with a strongly typed capability pointer.
func (m *Manager) DropGroupCapability2(name string, c *pb.Capability) error {
	if c == nil {
		return ErrUnknownCapability
	}

	rg := &pb.Group{
		Name:         &name,
		Capabilities: []pb.Capability{*c},
	}

	_, err := m.RunGroupChain("DROP-CAPABILITY", rg)
	return err
}
