package hooks

import (
	"context"

	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

// DestroyEntity removes an entity from the system.
type DestroyEntity struct {
	tree.BaseHook
}

// Run will request the underlying datastore to remove the entity,
// returning any status provided.  If the entity ID is not specified
// in e, it will be obtained from de.
func (d *DestroyEntity) Run(ctx context.Context, e, de *pb.Entity) error {
	// This hook is somewhat special since it might be called
	// after a processing pipeline, or just to remove an entity.
	if e.GetID() == "" {
		e.ID = de.ID
	}
	return d.Storage().DeleteEntity(ctx, e.GetID())
}

func init() {
	startup.RegisterCallback(destroyEntityCB)
}

func destroyEntityCB() {
	tree.RegisterEntityHookConstructor("destroy-entity", NewDestroyEntity)
}

// NewDestroyEntity returns an initialized DestroyEntity hook for use.
func NewDestroyEntity(opts ...tree.HookOption) (tree.EntityHook, error) {
	opts = append([]tree.HookOption{
		tree.WithHookName("destroy-entity"),
		tree.WithHookPriority(99),
	}, opts...)
	return &DestroyEntity{tree.NewBaseHook(opts...)}, nil
}
