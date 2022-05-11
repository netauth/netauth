package hooks

import (
	"context"

	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

// SaveEntity is designed to be a terminal processor in a chain.  On
// success, the provided entity will be saved to the data store.
type SaveEntity struct {
	tree.BaseHook
}

// Run will pass e to the data storage mechanism's "SaveEntity"
// method.
func (s *SaveEntity) Run(ctx context.Context, e, de *pb.Entity) error {
	return s.Storage().SaveEntity(ctx, e)
}

func init() {
	startup.RegisterCallback(saveEntityCB)
}

func saveEntityCB() {
	tree.RegisterEntityHookConstructor("save-entity", NewSaveEntity)
}

// NewSaveEntity returns an initialized hook ready for use.
func NewSaveEntity(opts ...tree.HookOption) (tree.EntityHook, error) {
	opts = append([]tree.HookOption{
		tree.WithHookName("save-entity"),
		tree.WithHookPriority(99),
	}, opts...)

	return &SaveEntity{tree.NewBaseHook(opts...)}, nil
}
