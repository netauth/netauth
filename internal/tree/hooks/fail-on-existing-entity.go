package hooks

import (
	"github.com/NetAuth/NetAuth/internal/db"
	"github.com/NetAuth/NetAuth/internal/tree/errors"

	pb "github.com/NetAuth/Protocol"
)

type FailOnExistingEntity struct {
	db.DB
}

func (*FailOnExistingEntity) Name() string  { return "fail-on-existing-entity" }
func (*FailOnExistingEntity) Priority() int { return 0 }
func (l *FailOnExistingEntity) Run(e, de *pb.Entity) error {
	_, err := l.LoadEntity(de.GetID())
	if err == nil {
		return tree.ErrDuplicateEntityID
	}
	return nil
}

