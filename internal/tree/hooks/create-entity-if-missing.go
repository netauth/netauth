package hooks

import (
	"github.com/NetAuth/NetAuth/internal/db"
	"github.com/NetAuth/NetAuth/internal/tree"
	"github.com/golang/protobuf/proto"

	pb "github.com/NetAuth/Protocol"
)

type CreateEntityIfMissing struct {
	tree.BaseHook
	db.DB
}

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

func init() {
	tree.RegisterEntityHookConstructor("create-entity-if-missing", NewCreateEntityIfMissing)
}

func NewCreateEntityIfMissing(c tree.RefContext) (tree.EntityProcessorHook, error) {
	return &CreateEntityIfMissing{tree.NewBaseHook("create-entity-if-missing", 1), c.DB}, nil
}
