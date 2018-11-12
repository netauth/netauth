package hooks

import (
	"github.com/NetAuth/NetAuth/internal/db"
	"github.com/NetAuth/NetAuth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

// SaveEntity is designed to be a terminal processor in a chain.  On
// success, the provided entity will be saved to the data store.
type SaveEntity struct {
	tree.BaseHook
	db.DB
}

// Run will pass e to the data storage mechanism's "SaveEntity"
// method.
func (s *SaveEntity) Run(e, de *pb.Entity) error {
	return s.SaveEntity(e)
}

func init() {
	tree.RegisterEntityHookConstructor("save-entity", NewSaveEntity)
}

// NewSaveEntity returns an initialized hook ready for use.
func NewSaveEntity(c tree.RefContext) (tree.EntityProcessorHook, error) {
	return &SaveEntity{tree.NewBaseHook("save-entity", 99), c.DB}, nil
}
