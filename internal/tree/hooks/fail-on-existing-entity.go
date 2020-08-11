package hooks

import (
	"github.com/netauth/netauth/internal/db"
	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

// FailOnExistingEntity checks with the data store to see if an entity
// exists.  If one does, then the hook will return an error that a
// duplicate ID already exists.
type FailOnExistingEntity struct {
	tree.BaseHook
	db.DB
}

// Run contacts the data store, attempts to load an entity and
// selectively inverts the return status from the load call (errors
// from the storage backend will be returned to the caller).
func (l *FailOnExistingEntity) Run(e, de *pb.Entity) error {
	_, err := l.LoadEntity(de.GetID())
	if err == nil {
		return tree.ErrDuplicateEntityID
	}
	return nil
}

func init() {
	startup.RegisterCallback(failOnExistingEntityCB)
}

func failOnExistingEntityCB() {
	tree.RegisterEntityHookConstructor("fail-on-existing-entity", NewFailOnExistingEntity)
}

// NewFailOnExistingEntity will return an initialized hook ready for
// use.
func NewFailOnExistingEntity(c tree.RefContext) (tree.EntityHook, error) {
	return &FailOnExistingEntity{tree.NewBaseHook("fail-on-existing-entity", 0), c.DB}, nil
}
