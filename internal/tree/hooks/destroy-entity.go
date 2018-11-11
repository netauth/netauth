package hooks

import (
	"github.com/NetAuth/NetAuth/internal/db"
	"github.com/NetAuth/NetAuth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

type DestroyEntity struct {
	tree.BaseHook
	db.DB
}

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

func NewDestroyEntity(c tree.RefContext) (tree.EntityProcessorHook, error) {
	return &DestroyEntity{tree.NewBaseHook("destroy-entity", 99), c.DB}, nil
}
