package hooks

import (
	pb "github.com/NetAuth/Protocol"
)

type SetEntityID struct{}

func (*SetEntityID) Name() string  { return "set-entity-ID" }
func (*SetEntityID) Priority() int { return 10 }
func (*SetEntityID) Run(e, de *pb.Entity) error {
	e.ID = de.ID
	return nil
}
