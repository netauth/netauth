package hooks

import (
	"github.com/netauth/netauth/internal/db"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

// DestroyEntity removes an entity from the system.
type DestroyEntity struct {
	tree.BaseHook
	db.DB
}

// Run will request the underlying datastore to remove the entity,
// returning any status provided.  If the entity ID is not specified
// in e, it will be obtained from de.
func (d *DestroyEntity) Run(e, de *pb.Entity) error {
	// This hook is somewhat special since it might be called
	// after a processing pipeline, or just to remove an entity.
	if e.GetID() == "" {
		e.ID = de.ID
	}
	return d.DeleteEntity(e.GetID())
}

func init() {
	tree.RegisterEntityHookConstructor("destroy-entity", NewDestroyEntity)
}

// NewDestroyEntity returns an initialized DestroyEntity hook for use.
func NewDestroyEntity(c tree.RefContext) (tree.EntityHook, error) {
	return &DestroyEntity{tree.NewBaseHook("destroy-entity", 99), c.DB}, nil
}
