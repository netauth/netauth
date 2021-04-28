package mresolver

import (
	"github.com/the-maldridge/bsfilter"
)

// SyncDirectGroups updates the list of groups in the resolver for a
// given entity with whatever the list actually is now.
func (mr *MResolver) SyncDirectGroups(entity string, groups []string) {
	list := make(map[string]struct{}, len(groups))
	for i := range groups {
		list[groups[i]] = struct{}{}
	}
	mr.uMutex.Lock()
	mr.atom.dm[entity] = list
	mr.uMutex.Unlock()
	mr.l.Trace("Synced direct groups", "entity", entity, "groups", groups)
}

// RemoveEntity removes an entity from the map, this is meant to
// handle deletions of entities.
func (mr *MResolver) RemoveEntity(entity string) {
	mr.uMutex.Lock()
	delete(mr.atom.dm, entity)
	mr.uMutex.Unlock()
}

// SyncGroup provides the resolver with current infomation about a
// given group.  Information here strictly overwrites other
// information in the system, and may trigger a cascading membership
// recalculation.
func (mr *MResolver) SyncGroup(group string, include, exclude []string) {
	g := resolvableGroup{
		self:    group,
		include: include,
		exclude: exclude,
	}

	mr.gMutex.Lock()
	mr.atom.gc[group] = &g
	mr.gMutex.Unlock()

	// Need to resolve the group here, and re-resolve any
	// dependant groups that this one would have affected.
	mr.l.Trace("Syncing group", "group", group)
	mr.Resolve(group)

	// Log all groups that affect this one to the group affectors
	// set.  This is a map to all groups that if changed need to
	// re-resolve this one.
	addAffector := func(affectee, affector string) {
		if mr.atom.ga[affectee] == nil {
			mr.atom.ga[affectee] = make(map[string]struct{})
		}
		mr.atom.ga[affectee][affector] = struct{}{}
	}
	mr.gMutex.Lock()
	for _, g := range include {
		addAffector(g, group)
	}
	for _, g := range exclude {
		addAffector(g, group)
	}
	mr.gMutex.Unlock()

	// Propagate changes to all groups affected by this group
	mr.resolveChanges(group)
}

// Recursively propagate resolutions all the way down.
func (mr *MResolver) resolveChanges(group string) {
	mr.Resolve(group)
	mr.gMutex.RLock()
	cascade := mr.atom.ga[group]
	mr.gMutex.RUnlock()
	for g := range cascade {
		mr.resolveChanges(g)
	}
}

// RemoveGroup removes all references to a group in the relations
// maps.  It does not remove any references in the entity direct
// memberships as this is assumed to have been done prior to deleting
// the group.
func (mr *MResolver) RemoveGroup(group string) {
	mr.gMutex.Lock()
	delete(mr.atom.gc, group)
	delete(mr.atom.gr, group)
	delete(mr.atom.gt, group)
	delete(mr.atom.ga, group)

	mr.atom.gs.Del(group)

	for _, ga := range mr.atom.ga {
		delete(ga, group)
	}
	mr.gMutex.Unlock()
}

// Resolve flattens out the membership tree and associated subtrees
// for any group into a boolean expression of arbitrary complexity.
// This expression can later be used to determine if an entity
// contains a group, or if a group contains a given set of entities.
func (mr *MResolver) Resolve(group string) error {
	mr.gMutex.RLock()
	g, ok := mr.atom.gc[group]
	mr.gMutex.RUnlock()
	if !ok {
		mr.l.Error("Insufficient knowledge to resolve group", "group", group)
		return ErrInsufficientKnowledge
	}

	mr.l.Trace("Resolving Group", "group", group)

	res := []bsfilter.Symbol{}
	res = append(res, bsfilter.Symbol{T: bsfilter.SymbolLParen})
	res = append(res, bsfilter.Symbol{T: bsfilter.SymbolIdent, Ident: g.self})

	// These two loops can be collapsed by iterating over a map at
	// the higher level and swapping between a logical operator
	// between loop iterations.
	for _, gn := range g.include {
		res = append(res, bsfilter.Symbol{T: bsfilter.SymbolBinaryOr})
		res = append(res, bsfilter.Symbol{T: bsfilter.SymbolLParen})
		mr.gMutex.RLock()
		// Check if an existing resolution is known for this
		// group name.  If not, resolve that group and then
		// fetch the expansion back out.
		r, ok := mr.atom.gt[gn]
		mr.gMutex.RUnlock()
		if !ok {
			mr.l.Trace("Resolution not cached; recursing", "group", gn)
			// We have to unlock the mutex to allow the
			// lower invocation to re-aquire.  This can in
			// potentially deadlock if another thread is
			// reading and trying to resolve.
			if err := mr.Resolve(gn); err != nil {
				return err
			}
			mr.gMutex.RLock()
			// At this point this is populated from the
			// above call and can be retreived
			r = mr.atom.gt[gn]
			mr.gMutex.RUnlock()
		}
		res = append(res, r...)
		res = append(res, bsfilter.Symbol{T: bsfilter.SymbolRParen})
	}
	res = append(res, bsfilter.Symbol{T: bsfilter.SymbolRParen})
	for _, gn := range g.exclude {
		res = append(res, bsfilter.Symbol{T: bsfilter.SymbolBinaryAnd})
		res = append(res, bsfilter.Symbol{T: bsfilter.SymbolUnaryNot})
		res = append(res, bsfilter.Symbol{T: bsfilter.SymbolLParen})
		mr.gMutex.RLock()
		// Check if an existing resolution is known for this
		// group name.  If not, resolve that group and then
		// fetch the expansion back out.
		r, ok := mr.atom.gt[gn]
		mr.gMutex.RUnlock()
		if !ok {
			mr.l.Trace("Resolution not cached; recursing", "group", gn)
			// We have to unlock the mutex to allow the
			// lower invocation to re-aquire.  This can in
			// potentially deadlock if another thread is
			// reading and trying to resolve.
			if err := mr.Resolve(gn); err != nil {
				return err
			}
			mr.gMutex.RLock()
			// At this point this is populated from the
			// above call and can be retreived
			r = mr.atom.gt[gn]
			mr.gMutex.RUnlock()
		}
		res = append(res, r...)
		res = append(res, bsfilter.Symbol{T: bsfilter.SymbolRParen})
	}
	expr := bsfilter.NewFromTokens(res)
	mr.gMutex.Lock()
	mr.atom.gt[group] = res
	mr.atom.gr[group] = expr
	mr.atom.gs.Add(group, expr)
	mr.gMutex.Unlock()
	mr.l.Debug("Group Resolved", "group", group, "expression", expr)
	return nil
}

// MembersOfGroup returns a list of all entities that are a member of
// the specified group.
func (mr *MResolver) MembersOfGroup(group string) []string {
	mr.gMutex.RLock()
	exp, ok := mr.atom.gr[group]
	mr.gMutex.RUnlock()
	if !ok {
		return []string{}
	}
	return exp.FilterValues(mr.atom.dm)
}

// GroupsForEntity returns a string slice of groups that include a
// given entity.
func (mr *MResolver) GroupsForEntity(entity string) []string {
	mr.uMutex.RLock()
	vset, ok := mr.atom.dm[entity]
	mr.uMutex.RUnlock()
	if !ok {
		return []string{}
	}
	return mr.atom.gs.Filter(vset)
}
