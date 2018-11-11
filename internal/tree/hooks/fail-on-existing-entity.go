package hooks

import (
	"github.com/NetAuth/NetAuth/internal/db"
	"github.com/NetAuth/NetAuth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

type FailOnExistingEntity struct {
	tree.BaseHook
	db.DB
}

func (l *FailOnExistingEntity) Run(e, de *pb.Entity) error {
	_, err := l.LoadEntity(de.GetID())
	if err == nil {
		return tree.ErrDuplicateEntityID
	}
	return nil
}

func init() {
	tree.RegisterEntityHookConstructor("fail-on-existing-entity", NewFailOnExistingEntity)
}

func NewFailOnExistingEntity(c tree.RefContext) (tree.EntityProcessorHook, error) {
	return &FailOnExistingEntity{tree.NewBaseHook("fail-on-existing-entity", 0), c.DB}, nil
}
