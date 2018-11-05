package hooks

import (
	"github.com/NetAuth/NetAuth/internal/db"
	"github.com/golang/protobuf/proto"

	pb "github.com/NetAuth/Protocol"
)

type CreateEntityIfMissing struct {
	db.DB
}

func (*CreateEntityIfMissing) Name() string  { return "create-entity-if-missing" }
func (*CreateEntityIfMissing) Priority() int { return 10 }
func (c *CreateEntityIfMissing) Run(e, de *pb.Entity) error {
	le, err := c.LoadEntity(de.GetID())
	switch err {
	case nil:
		proto.Merge(e, le)
		return err
	case db.ErrUnknownEntity:
		break
	default:
		return err
	}

	ce := &pb.Entity{
		ID: de.ID,
	}
	proto.Merge(e, ce)
	return nil
}
