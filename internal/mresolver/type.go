package mresolver

import (
	"sync"

	"github.com/hashicorp/go-hclog"

	"github.com/the-maldridge/bsfilter"
)

// MResolver contains flattened data structures that resolve
// memberships for entities and which entities are members of a given
// group.
type MResolver struct {
	l hclog.Logger

	uMutex sync.RWMutex
	gMutex sync.RWMutex

	atom resolverAtom
}

type resolverAtom struct {
	dm map[string]bsfilter.ValueSet    // Cache of direct memberships
	gc map[string]*resolvableGroup     // Cache of groups and rules
	gr map[string]*bsfilter.Expression // Resolved expressions
	gt map[string][]bsfilter.Symbol    // Cache of subexpressions
	ga map[string]map[string]struct{}  // Cache of groups that are affected by the key group
	gs *bsfilter.ExpressionSet         // Set of all expressions for all groups
}

type resolvableGroup struct {
	self string

	include []string
	exclude []string
}
