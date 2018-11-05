package hooks

import (
	"github.com/NetAuth/NetAuth/internal/db"
	"github.com/golang/protobuf/proto"

	pb "github.com/NetAuth/Protocol"
)

type LoadEntity struct {
	db.DB
}

func (*LoadEntity) Name() string  { return "load-entity" }
func (*LoadEntity) Priority() int { return 5 }
func (l *LoadEntity) Run(e, de *pb.Entity) error {
	// This is a bit odd because we only get an address for e, not
	// the ability to point it somewhere else, so anything we want
	// to do that alters the initial contents needs to be copied
	// in.

	le, err := l.LoadEntity(de.GetID())
	if err != nil {
		return err
	}
	proto.Merge(e, le)

	return nil
}
