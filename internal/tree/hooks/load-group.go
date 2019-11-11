package hooks

import (
	"github.com/golang/protobuf/proto"
	"github.com/netauth/netauth/internal/db"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

// LoadGroup loads an entity from the database.
type LoadGroup struct {
	tree.BaseHook
	db.DB
}

// Run attempts to load the group specified by de.Name and if
// successful performs a deepcopy into the address pointed to by g.
// Any errors returned will be from the data storage layer.
func (l *LoadGroup) Run(g, dg *pb.Group) error {
	// This is a bit odd because we only get an address for g, not
	// the ability to point it somewhere else, so anything we want
	// to do that alters the initial contents needs to be copied
	// in.

	lg, err := l.LoadGroup(dg.GetName())
	if err != nil {
		return err
	}
	proto.Merge(g, lg)

	return nil
}

func init() {
	tree.RegisterGroupHookConstructor("load-group", NewLoadGroup)
}

// NewLoadGroup returns an initialized hook ready for use.
func NewLoadGroup(c tree.RefContext) (tree.GroupHook, error) {
	return &LoadGroup{tree.NewBaseHook("load-group", 0), c.DB}, nil
}
