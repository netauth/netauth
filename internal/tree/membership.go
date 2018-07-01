package tree

import (
	"fmt"
	"log"
	"strings"

	pb "github.com/NetAuth/Protocol"
)

// AddEntityToGroup is the same as the internal function, but takes an
// entity ID rather than a pointer
func (m *Manager) AddEntityToGroup(entityID, groupName string) error {
	e, err := m.db.LoadEntity(entityID)
	if err != nil {
		return err
	}
	return m.addEntityToGroup(e, groupName)
}

// addEntityToGroup adds an entity to a group by name, if the entity
// was already in the group the function will return with a nil error.
func (m *Manager) addEntityToGroup(e *pb.Entity, groupName string) error {
	if _, err := m.db.LoadGroup(groupName); err != nil {
		return err
	}

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

	if err := m.db.SaveEntity(e); err != nil {
		return err
	}
	return nil
}

// GetMemberships returns all groups the entity is a member of,
// optionally including indirect memberships
func (m *Manager) GetMemberships(e *pb.Entity, includeIndirects bool) []string {
	directs := m.getDirectGroups(e)
	var allGroups []string

	// Though inneficient, its easier to understand.  We get the
	// membership of all groups, and evaluate if this entity is a
	// member of those groups.
	grps, err := m.ListGroups()
	if err != nil {
		log.Printf("Error getting group list: %s", grps)
		return []string{}
	}
	for _, g := range grps {
		members, err := m.listMembers(g.GetName())
		if err != nil {
			log.Printf("Error expanding groups: %s", err)
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

// RemoveEntityFromGroup performs the same function as the internal
// variant, but does so by name rather than by entity pointer.
func (m *Manager) RemoveEntityFromGroup(entityID, groupName string) error {
	e, err := m.db.LoadEntity(entityID)
	if err != nil {
		return err
	}
	m.removeEntityFromGroup(e, groupName)
	return nil
}

// removeEntityFromGroup removes an entity from the named group.  If
// the entity was not in the group to begin with then nil will be
// returned as the error.
func (m *Manager) removeEntityFromGroup(e *pb.Entity, groupName string) error {
	if e.GetMeta() == nil {
		return nil
	}

	newGroups := []string{}
	for _, g := range e.GetMeta().GetGroups() {
		if g == groupName {
			continue
		}
		newGroups = append(newGroups, g)
	}
	e.Meta.Groups = newGroups

	if err := m.db.SaveEntity(e); err != nil {
		return err
	}
	return nil
}

// allEntities is a convenient way to return all the entities
func (m *Manager) allEntities() ([]*pb.Entity, error) {
	var entities []*pb.Entity
	el, err := m.db.DiscoverEntityIDs()
	if err != nil {
		return nil, err
	}
	for _, en := range el {
		e, err := m.db.LoadEntity(en)
		if err != nil {
			return nil, err
		}
		entities = append(entities, e)
	}
	return entities, nil
}

// listMembers takes a group ID in and returns a slice of entities
// that are in that group.
func (m *Manager) listMembers(groupID string) ([]*pb.Entity, error) {
	// 'ALL' is a special groupID which returns everything, this
	// isn't a group that exists in a real sense, it just serves
	// to return a global list as a convenience.
	if groupID == "ALL" {
		return m.allEntities()
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
	el, err := m.allEntities()
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
			log.Printf("Expansion parsing error! %s", err)
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

// checkExistingGroupExpansions verifies that there is no expansion
// already directly on this group that conflicts with the proposed
// group expansion.
func (m *Manager) checkExistingGroupExpansions(g *pb.Group, candidate string) error {
	for _, exp := range g.GetExpansions() {
		if strings.Contains(exp, candidate) {
			return ErrExistingExpansion
		}
	}
	return nil
}

// checkGroupCycles recurses down the group tree and tries to find the
// candidate group somewhere on the tree below the entry point.  The
// general usage would be to push in the target of the expansion as
// the group and then hunt for the parent group as the candidate.
func (m *Manager) checkGroupCycles(g *pb.Group, candidate string) bool {
	for _, exp := range g.GetExpansions() {
		parts := strings.Split(exp, ":")
		log.Println(parts[1], candidate)
		if parts[1] == candidate {
			return true
		}
		g, err := m.GetGroupByName(parts[1])
		if err != nil {
			// Play it safe, if we can't get the group
			// something may already be wrong.  Returning
			// true here can prevent further damage to the
			// tree.
			return true
		}
		return m.checkGroupCycles(g, candidate)
	}
	return false
}

// ModifyGroupExpansions handles changing the expansions on a group.
// This can include adding an INCLUDE or EXCLUDE type expansion, or
// using the special expansion type DROP, removing an existing one.
func (m *Manager) ModifyGroupExpansions(parent, child string, mode pb.ExpansionMode) error {
	p, err := m.GetGroupByName(parent)
	if err != nil {
		return err
	}

	// Check if there are any conflicting direct expansions on
	// this group.  Expansions on children are fine if they
	// conflict, that will just be confusing, but a conflicting
	// expansion here could cause undefined behavior.
	if err := m.checkExistingGroupExpansions(p, child); err != nil && mode != pb.ExpansionMode_DROP {
		return err
	}

	// Make sure the child exists...
	c, err := m.GetGroupByName(child)
	if err != nil {
		return err
	}

	if m.checkGroupCycles(c, p.GetName()) && mode != pb.ExpansionMode_DROP {
		return ErrExistingExpansion
	}

	// Either add the include, add the exclude, or drop the old
	// record.
	switch mode {
	case pb.ExpansionMode_INCLUDE:
		p.Expansions = append(p.Expansions, fmt.Sprintf("%s:%s", mode, c.GetName()))
	case pb.ExpansionMode_EXCLUDE:
		p.Expansions = append(p.Expansions, fmt.Sprintf("%s:%s", mode, c.GetName()))
	case pb.ExpansionMode_DROP:
		old := p.GetExpansions()
		new := []string{}
		for _, oldMembership := range old {
			if strings.Contains(oldMembership, child) {
				continue
			}
			new = append(new, oldMembership)
		}
		p.Expansions = new
	}

	return m.db.SaveGroup(p)
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
