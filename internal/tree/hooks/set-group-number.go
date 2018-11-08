package hooks

import (
	"github.com/NetAuth/NetAuth/internal/db"

	pb "github.com/NetAuth/Protocol"
)

type SetGroupNumber struct {
	db.DB
}

func (*SetGroupNumber) Name() string { return "set-group-number" }
func (*SetGroupNumber) Priority() int { return 50 }
func (s *SetGroupNumber) Run(g, dg *pb.Group) error {
	if dg.GetNumber() == -1 {
		number, err := s.NextGroupNumber()
		if err != nil {
			return err
		}
		g.Number = &number
		return nil
	}
	dg.Number = g.Number
	return nil
}
