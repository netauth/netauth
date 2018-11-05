package hooks

import (
	"github.com/NetAuth/NetAuth/internal/db"

	pb "github.com/NetAuth/Protocol"
)

type DestroyEntity struct {
	db.DB
}

func (*DestroyEntity) Name() string  { return "destroy-entity" }
func (*DestroyEntity) Priority() int { return 99 }
func (d *DestroyEntity) Run(e, de *pb.Entity) error {
	// This hook is somewhat special since it might be called
	// after a processing pipeline, or just to remove an entity.
	if e.GetID() == "" {
		e.ID = de.ID
	}
	return d.DeleteEntity(e.GetID())
}
