package hooks

import (
	"github.com/NetAuth/NetAuth/internal/db"
	"github.com/NetAuth/NetAuth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

type SetEntityNumber struct {
	tree.BaseHook
	db.DB
}

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

func init() {
	tree.RegisterEntityHookConstructor("set-entity-number", NewSetEntityNumber)
}

func NewSetEntityNumber(c tree.RefContext) (tree.EntityProcessorHook, error) {
	return &SetEntityNumber{tree.NewBaseHook("set-entity-number", 50), c.DB}, nil
}
