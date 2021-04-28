package mresolver

import (
	"github.com/hashicorp/go-hclog"

	"github.com/the-maldridge/bsfilter"
)

// New sets up a new resolver.
func New() *MResolver {
	return &MResolver{
		l: hclog.NewNullLogger(),
		atom: resolverAtom{
			dm: make(map[string]bsfilter.ValueSet),
			gc: make(map[string]*resolvableGroup),
			gr: make(map[string]*bsfilter.Expression),
			gt: make(map[string][]bsfilter.Symbol),
			ga: make(map[string]map[string]struct{}),
			gs: bsfilter.NewExpressionSet(),
		},
	}
}

// SetParentLogger provides the parent logger from which the resolver
// will derive a private logger.
func (mr *MResolver) SetParentLogger(l hclog.Logger) {
	mr.l = l.Named("resolver")
}
