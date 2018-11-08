package hooks

import (
	"github.com/NetAuth/NetAuth/internal/db"
	"github.com/NetAuth/NetAuth/internal/tree/errors"

	pb "github.com/NetAuth/Protocol"
)

type FailOnExistingGroup struct {
	db.DB
}

func (*FailOnExistingGroup) Name() string { return "fail-on-existing-group" }
func (*FailOnExistingGroup) Priority() int { return 0 }
func (f *FailOnExistingGroup) Run(g, dg *pb.Group) error {
	if _, err := f.LoadGroup(dg.GetName()); err == nil {
		return tree.ErrDuplicateGroupName
	} else if err != db.ErrUnknownGroup {
		return err
	}
	return nil
}
