package hooks

import (
	"github.com/golang/protobuf/proto"

	"github.com/netauth/netauth/internal/db"
	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

// LoadEntity loads an entity from the database.
type LoadEntity struct {
	tree.BaseHook
	db.DB
}

// Run attempts to load the entity specified by de.ID and if
// successful performs a deepcopy into the address pointed to by e.
// Any errors returned will be from the data storage layer.
func (l *LoadEntity) Run(e, de *pb.Entity) error {
	// This is a bit odd because we only get an address for e, not
	// the ability to point it somewhere else, so anything we want
	// to do that alters the initial contents needs to be copied
	// in.

	le, err := l.LoadEntity(de.GetID())
	if err != nil {
		return err
	}
	proto.Merge(e, le)

	return nil
}

func init() {
	startup.RegisterCallback(loadEntityCB)
}

func loadEntityCB() {
	tree.RegisterEntityHookConstructor("load-entity", NewLoadEntity)
}

// NewLoadEntity returns an initialized hook ready for use.
func NewLoadEntity(c tree.RefContext) (tree.EntityHook, error) {
	return &LoadEntity{tree.NewBaseHook("load-entity", 0), c.DB}, nil
}
