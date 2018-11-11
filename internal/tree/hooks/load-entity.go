package hooks

import (
	"github.com/NetAuth/NetAuth/internal/db"
	"github.com/NetAuth/NetAuth/internal/tree"
	"github.com/golang/protobuf/proto"

	pb "github.com/NetAuth/Protocol"
)

type LoadEntity struct {
	tree.BaseHook
	db.DB
}

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

func init() {
	tree.RegisterEntityHookConstructor("load-entity", NewLoadEntity)
}

func NewLoadEntity(c tree.RefContext) (tree.EntityProcessorHook, error) {
	return &LoadEntity{tree.NewBaseHook("load-entity", 0), c.DB}, nil
}
