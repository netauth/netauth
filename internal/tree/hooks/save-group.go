package hooks

import (
	"github.com/NetAuth/NetAuth/internal/db"

	pb "github.com/NetAuth/Protocol"
)

type SaveGroup struct {
	db.DB
}

func (*SaveGroup) Name() string { return "save-group" }
func (*SaveGroup) Priority() int { return 99 }
func (s *SaveGroup) Run(g, dg *pb.Group) error {
	return s.SaveGroup(g)
}
