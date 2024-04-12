package hooks

import (
	"context"

	"github.com/netauth/netauth/internal/db"
	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

// FailOnEmptyEntity checks with the data store to see if an entity
// exists.  If one does, then the hook will return an error that a
// duplicate ID already exists.
type FailOnEmptyEntity struct {
	tree.BaseHook
}

// Run contacts the data store, attempts to load an entity and
// selectively inverts the return status from the load call (errors
// from the storage backend will be returned to the caller).
func (l *FailOnEmptyEntity) Run(ctx context.Context, e, de *pb.Entity) error {
	if de.GetID() == "" {
		return db.ErrNoValue
	}
	return nil
}

func init() {
	startup.RegisterCallback(failOnEmptyEntityCB)
}

func failOnEmptyEntityCB() {
	tree.RegisterEntityHookConstructor("fail-on-empty-entity", NewFailOnEmptyEntity)
}

// NewFailOnEmptyEntity will return an initialized hook ready for
// use.
func NewFailOnEmptyEntity(opts ...tree.HookOption) (tree.EntityHook, error) {
	opts = append([]tree.HookOption{
		tree.WithHookName("fail-on-empty-entity"),
		tree.WithHookPriority(0),
	}, opts...)

	return &FailOnEmptyEntity{tree.NewBaseHook(opts...)}, nil
}
