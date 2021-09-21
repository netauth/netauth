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
	tree.DB
}

// Run will pass e to the data storage mechanism's "SaveEntity"
// method.
func (s *SaveEntity) Run(ctx context.Context, e, de *pb.Entity) error {
	return s.SaveEntity(ctx, e)
}

func init() {
	startup.RegisterCallback(saveEntityCB)
}

func saveEntityCB() {
	tree.RegisterEntityHookConstructor("save-entity", NewSaveEntity)
}

// NewSaveEntity returns an initialized hook ready for use.
func NewSaveEntity(c tree.RefContext) (tree.EntityHook, error) {
	return &SaveEntity{tree.NewBaseHook("save-entity", 99), c.DB}, nil
}
