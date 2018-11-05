package hooks

import (
	"github.com/NetAuth/NetAuth/internal/db"

	pb "github.com/NetAuth/Protocol"
)

type SaveEntity struct {
	db.DB
}

func (*SaveEntity) Name() string  { return "save-entity" }
func (*SaveEntity) Priority() int { return 99 }
func (s *SaveEntity) Run(e, de *pb.Entity) error {
	return s.SaveEntity(e)
}
