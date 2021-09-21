package tree

import (
	"context"
	"fmt"

	pb "github.com/netauth/protocol"
	rpc "github.com/netauth/protocol/v2"
)

// AddEntityToGroup is the same as the internal function, but takes an
// entity ID rather than a pointer
func (m *Manager) AddEntityToGroup(ctx context.Context, entityID, groupName string) error {
	de := &pb.Entity{
		ID: &entityID,
		Meta: &pb.EntityMeta{
			Groups: []string{groupName},
		},
	}

	_, err := m.RunEntityChain(ctx, "GROUP-ADD", de)
	return err
}

// RemoveEntityFromGroup performs the same function as the internal
// variant, but does so by name rather than by entity pointer.
func (m *Manager) RemoveEntityFromGroup(ctx context.Context, entityID, groupName string) error {
	de := &pb.Entity{
		ID: &entityID,
		Meta: &pb.EntityMeta{
			Groups: []string{groupName},
		},
	}

	_, err := m.RunEntityChain(ctx, "GROUP-DEL", de)
	return err
}

// GetMemberships returns a list of group names that an entity is a
// member of.  This membership may either be direct or it may be via
// an expanded group rule.  This difference is not distinguished.
func (m *Manager) GetMemberships(ctx context.Context, e *pb.Entity) []string {
	return m.resolver.GroupsForEntity(e.GetID())
}

// ListMembers fetches the members of a single group and redacts
// authentication data.
func (m *Manager) ListMembers(ctx context.Context, groupID string) ([]*pb.Entity, error) {
	eIDs := m.resolver.MembersOfGroup(groupID)
	m.log.Trace("Resolved Entities", "group", groupID, "entities", eIDs)

	entities := []*pb.Entity{}
	for i := range eIDs {
		e, err := m.db.LoadEntity(ctx, eIDs[i])
		if err != nil {
			return nil, err
		}
		entities = append(entities, e)
	}

	var safeMembers []*pb.Entity
	for _, e := range entities {
		ne := safeCopyEntity(e)
		safeMembers = append(safeMembers, ne)
	}

	return safeMembers, nil
}

// ModifyGroupExpansions handles changing the expansions on a group.
// This can include adding an INCLUDE or EXCLUDE type expansion, or
// using the special expansion type DROP, removing an existing one.
func (m *Manager) ModifyGroupExpansions(ctx context.Context, parent, child string, mode pb.ExpansionMode) error {
	rg := &pb.Group{
		Name:       &parent,
		Expansions: []string{fmt.Sprintf("%s:%s", mode, child)},
	}
	_, err := m.RunGroupChain(ctx, "MODIFY-EXPANSIONS", rg)
	return err
}

// ModifyGroupRule adjusts the rules on a group, which is the second
// iteration of the expansion system.  Right now this function is a
// shim over the legacy ModifyGroupExpansions interface, but it will
// be modified to support the strongly typed group interface at a
// later date.
func (m *Manager) ModifyGroupRule(ctx context.Context, group, target string, ruleaction rpc.RuleAction) error {
	mode := "DROP"
	switch ruleaction {
	case rpc.RuleAction_INCLUDE:
		mode = "INCLUDE"
	case rpc.RuleAction_EXCLUDE:
		mode = "EXCLUDE"
	case rpc.RuleAction_REMOVE_RULE:
		mode = "DROP"
	}

	// We can do this un-checked since there's a hard coded list
	// of strings above that are used to select this value.
	md := pb.ExpansionMode(pb.ExpansionMode_value[mode])

	return m.ModifyGroupExpansions(ctx, group, target, md)
}
