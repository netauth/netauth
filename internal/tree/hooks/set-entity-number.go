package hooks

import (
	"github.com/NetAuth/NetAuth/internal/db"

	pb "github.com/NetAuth/Protocol"
)

type SetEntityNumber struct {
	db.DB
}

func (*SetEntityNumber) Name() string  { return "set-entity-number" }
func (*SetEntityNumber) Priority() int { return 5 }
func (s *SetEntityNumber) Run(e, de *pb.Entity) error {
	if de.GetNumber() == -1 {
		n, err := s.NextEntityNumber()
		if err != nil {
			return err
		}
		e.Number = &n
		return nil
	}
	e.Number = de.Number
	return nil
}

