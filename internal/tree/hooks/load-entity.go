package hooks

import (
	"github.com/NetAuth/NetAuth/internal/db"

	pb "github.com/NetAuth/Protocol"
)

type loadEntity struct {
	db.DB
}

func (*loadEntity) Name() string  { return "load-entity" }
func (*loadEntity) Priority() int { return 5 }
func (l *loadEntity) Run(e, de *pb.Entity) error {
	e, err := l.LoadEntity(de.GetID())
	if err != nil {
		return err
	}
	return nil
}
