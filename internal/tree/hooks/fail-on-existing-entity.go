package hooks

import (
	"context"

	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

// FailOnExistingEntity checks with the data store to see if an entity
// exists.  If one does, then the hook will return an error that a
// duplicate ID already exists.
type FailOnExistingEntity struct {
	tree.BaseHook
}

// Run contacts the data store, attempts to load an entity and
// selectively inverts the return status from the load call (errors
// from the storage backend will be returned to the caller).
func (l *FailOnExistingEntity) Run(ctx context.Context, e, de *pb.Entity) error {
	_, err := l.Storage().LoadEntity(ctx, de.GetID())
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
func NewFailOnExistingEntity(opts ...tree.HookOption) (tree.EntityHook, error) {
	opts = append([]tree.HookOption{
		tree.WithHookName("fail-on-existing-entity"),
		tree.WithHookPriority(0),
	}, opts...)

	return &FailOnExistingEntity{tree.NewBaseHook(opts...)}, nil
}
