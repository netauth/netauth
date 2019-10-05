package tree

import (
	"fmt"
	"strings"

	"github.com/NetAuth/NetAuth/internal/db"

	pb "github.com/NetAuth/Protocol"
	rpc "github.com/NetAuth/Protocol/v2"
)

// AddEntityToGroup is the same as the internal function, but takes an
// entity ID rather than a pointer
func (m *Manager) AddEntityToGroup(entityID, groupName string) error {
	de := &pb.Entity{
		ID: &entityID,
		Meta: &pb.EntityMeta{
			Groups: []string{groupName},
		},
	}

	_, err := m.RunEntityChain("GROUP-ADD", de)
	return err
}

// RemoveEntityFromGroup performs the same function as the internal
// variant, but does so by name rather than by entity pointer.
func (m *Manager) RemoveEntityFromGroup(entityID, groupName string) error {
	de := &pb.Entity{
		ID: &entityID,
		Meta: &pb.EntityMeta{
			Groups: []string{groupName},
		},
	}

	_, err := m.RunEntityChain("GROUP-DEL", de)
	return err
}

// GetMemberships returns all groups the entity is a member of,
// optionally including indirect memberships
func (m *Manager) GetMemberships(e *pb.Entity, includeIndirects bool) []string {
	directs := m.getDirectGroups(e)
	var allGroups []string

	// Though inneficient, its easier to understand.  We get the
	// membership of all groups, and evaluate if this entity is a
	// member of those groups.
	grps, err := m.SearchGroups(db.SearchRequest{Expression: "*"})
	if err != nil {
		m.log.Error("Error getting complete group list", "error", err)
		return []string{}
	}
	for _, g := range grps {
		members, err := m.listMembers(g.GetName())
		if err != nil {
			m.log.Warn("Error during group expansion", "group", g.GetName(), "error", err)
			continue
		}
		for _, m := range members {
			if m.GetID() == e.GetID() {
				allGroups = append(allGroups, g.GetName())
				break
			}
		}
	}

	// If we're including indirects, then we can return allGroups
	// here
	if includeIndirects {
		return allGroups
	}

	// This far?  Only returning directs as filtered by allGroups.
	// This is because there could be things that filter entities
	// out of groups they would otherwise be directly in.
	gm := make(map[string]int)
	for _, g := range directs {
		gm[g]++
	}
	for _, g := range allGroups {
		gm[g]++
	}

	// gm now contains a map where groups that are both in
	// allGroups (membership is valid) and in directs (groups we
	// want to return here) will have a value of 2, so all we do
	// at this stage is get the groups that are equal to 2 and
	// return those.
	var retGroups []string
	for name, value := range gm {
		if value == 2 {
			retGroups = append(retGroups, name)
		}
	}
	return retGroups
}

// getDirectGroups gets the direct groups of an entity.
func (m *Manager) getDirectGroups(e *pb.Entity) []string {
	if e.GetMeta() == nil {
		return []string{}
	}

	return e.GetMeta().GetGroups()
}

// listMembers takes a group ID in and returns a slice of entities
// that are in that group.
func (m *Manager) listMembers(groupID string) ([]*pb.Entity, error) {
	// 'ALL' is a special groupID which returns everything, this
	// isn't a group that exists in a real sense, it just serves
	// to return a global list as a convenience.
	if groupID == "ALL" {
		return m.SearchEntities(db.SearchRequest{Expression: "*"})
	}

	// If its not the all group then we check to make sure the
	// group exists at all
	g, err := m.db.LoadGroup(groupID)
	if err != nil {
		return nil, err
	}

	// Now we can be reasonably sure the group exists, this next
	// bit is stupidly inefficient, but is the only way to extract
	// the members since the membership graph has the arrows going
	// the other way.
	var entities []*pb.Entity
	el, err := m.SearchEntities(db.SearchRequest{Expression: "*"})
	if err != nil {
		return nil, err
	}
	for _, e := range el {
		for _, g := range m.getDirectGroups(e) {
			if g == groupID {
				entities = append(entities, e)
			}
		}
	}

	// Now we parse the expansions.
	var exclude []*pb.Entity
	for _, exp := range g.GetExpansions() {
		parts := strings.Split(exp, ":")
		ents, err := m.listMembers(parts[1])
		if err != nil {
			m.log.Error("Error parsing expansion", "expansion", exp, "error", err)
		}
		switch parts[0] {
		case "INCLUDE":
			entities = append(entities, ents...)
		case "EXCLUDE":
			exclude = append(exclude, ents...)
		}
	}

	// Its likely that we've got duplicates in the lists now, so
	// dedup things to get back down to one copy of everything.
	entities = dedupEntityList(entities)
	exclude = dedupEntityList(exclude)

	// Actually exclude the excluded entities
	if len(exclude) > 0 {
		entities = entityListDifference(entities, exclude)
	}

	// This will be the entities that are in this group and all of
	// its expansions, but not any that would be excluded from
	// this group or the subexpansions.
	return entities, nil
}

// ListMembers fulfills the same function as the private version of
// this function, but with one crucial difference, it produces copies
// of the entities that have the secret redacted.
func (m *Manager) ListMembers(groupID string) ([]*pb.Entity, error) {
	// This set of members has secrets and can't be returned.
	members, err := m.listMembers(groupID)
	if err != nil {
		return nil, err
	}

	var safeMembers []*pb.Entity
	for _, e := range members {
		ne := safeCopyEntity(e)
		safeMembers = append(safeMembers, ne)
	}

	return safeMembers, nil
}

// ModifyGroupExpansions handles changing the expansions on a group.
// This can include adding an INCLUDE or EXCLUDE type expansion, or
// using the special expansion type DROP, removing an existing one.
func (m *Manager) ModifyGroupExpansions(parent, child string, mode pb.ExpansionMode) error {
	rg := &pb.Group{
		Name:       &parent,
		Expansions: []string{fmt.Sprintf("%s:%s", mode, child)},
	}
	_, err := m.RunGroupChain("MODIFY-EXPANSIONS", rg)
	return err
}

// ModifyGroupRule adjusts the rules on a group, which is the second
// iteration of the expansion system.  Right now this function is a
// shim over the legacy ModifyGroupExpansions interface, but it will
// be modified to support the strongly typed group interface at a
// later date.
func (m *Manager) ModifyGroupRule(group, target string, ruleaction rpc.RuleAction) error {
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

	return m.ModifyGroupExpansions(group, target, md)
}

// dedupEntityList takes in a list of entities and deduplicates them
// using a map.
func dedupEntityList(entList []*pb.Entity) []*pb.Entity {
	eMap := make(map[string]*pb.Entity)
	for _, e := range entList {
		eMap[e.GetID()] = e
	}

	// Back to a list...
	var eList []*pb.Entity
	for _, e := range eMap {
		eList = append(eList, e)
	}
	return eList
}

// entityListDifference computes the set of entities that are in list
// a and not in list b.
func entityListDifference(a, b []*pb.Entity) []*pb.Entity {
	diffMap := make(map[string]*pb.Entity)
	// Get a map of the possible options
	for _, e := range a {
		diffMap[e.GetID()] = e
	}
	// Remove the ones that are in the exclude map
	for _, e := range b {
		delete(diffMap, e.GetID())
	}

	var entList []*pb.Entity
	for _, ent := range diffMap {
		entList = append(entList, ent)
	}

	return entList
}
