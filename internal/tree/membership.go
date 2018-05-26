package tree

import (
	"fmt"
	"strings"

	"github.com/NetAuth/NetAuth/pkg/errors"
	pb "github.com/NetAuth/Protocol"
)

// AddEntityToGroup is the same as the internal function, but takes an
// entity ID rather than a pointer
func (m Manager) AddEntityToGroup(entityID, groupName string) error {
	e, err := m.db.LoadEntity(entityID)
	if err != nil {
		return err
	}
	return m.addEntityToGroup(e, groupName)
}

// addEntityToGroup adds an entity to a group by name, if the entity
// was already in the group the function will return with a nil error.
func (m Manager) addEntityToGroup(e *pb.Entity, groupName string) error {
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
func (m Manager) GetMemberships(e *pb.Entity, includeIndirects bool) []string {
	return m.GetDirectGroups(e)
}

// GetDirectGroups gets the direct groups of an entity.
func (m Manager) GetDirectGroups(e *pb.Entity) []string {
	if e.GetMeta() == nil {
		return []string{}
	}

	return e.GetMeta().GetGroups()
}

// RemoveEntityFromGroup performs the same function as the internal
// variant, but does so by name rather than by entity pointer.
func (m Manager) RemoveEntityFromGroup(entityID, groupName string) error {
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
func (m Manager) removeEntityFromGroup(e *pb.Entity, groupName string) {
	if e.GetMeta() == nil {
		return
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
		return
	}
	return
}

// listMembers takes a group ID in and returns a slice of entities
// that are in that group.
func (m Manager) listMembers(groupID string) ([]*pb.Entity, error) {
	var entities []*pb.Entity

	// 'ALL' is a special groupID which returns everything, this
	// isn't a group that exists in a real sense, it just serves
	// to return a global list as a convenience.
	if groupID == "ALL" {
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

	// If its not the all group then we check to make sure the
	// group exists at all
	if _, err := m.db.LoadGroup(groupID); err != nil {
		return nil, err
	}

	// Now we can be reasonably sure the group exists, this next
	// bit is stupidly inefficient, but is the only way to extract
	// the members since the membership graph has the arrows going
	// the other way.
	el, err := m.listMembers("ALL")
	if err != nil {
		return nil, err
	}
	for _, e := range el {
		for _, g := range m.GetDirectGroups(e) {
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
func (m Manager) ListMembers(groupID string) ([]*pb.Entity, error) {
	// This set of members has secrets and can't be returned.
	members, err := m.listMembers(groupID)
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

// checkExistingGroupExpansions verifies that there is no expansion
// already directly on this group that conflicts with the proposed
// group expansion.
func (m Manager) checkExistingGroupExpansions(g *pb.Group, candidate string) error {
	for _, exp := range g.GetChildren() {
		if strings.Contains(exp, candidate) {
			return errors.E_EXISTING_EXPANSION
		}
	}
	return nil
}

// ModifyGroupExpansions handles changing the expansions on a group.
// This can include adding an INCLUDE or EXCLUDE type expansion, or
// using the special expansion type DROP, removing an existing one.
func (m Manager) ModifyGroupExpansions(parent, child string, mode pb.ExpansionMode) error {
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

	// Either add the include, add the exclude, or drop the old
	// record.
	switch mode {
	case pb.ExpansionMode_INCLUDE:
		p.Children = append(p.Children, fmt.Sprintf("%s:%s", mode, c.GetName()))
	case pb.ExpansionMode_EXCLUDE:
		p.Children = append(p.Children, fmt.Sprintf("%s:%s", mode, c.GetName()))
	case pb.ExpansionMode_DROP:
		old := p.GetChildren()
		new := []string{}
		for _, oldMembership := range old {
			if strings.Contains(oldMembership, child) {
				continue
			}
			new = append(new, oldMembership)
		}
		p.Children = new
	}

	return m.db.SaveGroup(p)
}
