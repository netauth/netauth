package hooks

import (
	"github.com/NetAuth/NetAuth/internal/db"
	"github.com/NetAuth/NetAuth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

type SaveEntity struct {
	tree.BaseHook
	db.DB
}

func (s *SaveEntity) Run(e, de *pb.Entity) error {
	return s.SaveEntity(e)
}

func init() {
	tree.RegisterEntityHookConstructor("save-entity", NewSaveEntity)
}

func NewSaveEntity(c tree.RefContext) (tree.EntityProcessorHook, error) {
	return &SaveEntity{tree.NewBaseHook("save-entity", 99), c.DB}, nil
}
