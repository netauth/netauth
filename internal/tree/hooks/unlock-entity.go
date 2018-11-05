package hooks

import (
	"github.com/golang/protobuf/proto"

	pb "github.com/NetAuth/Protocol"
)

type UnlockEntity struct{}

func (*UnlockEntity) Name() string  { return "unlock-entity" }
func (*UnlockEntity) Priority() int { return 50 }
func (*UnlockEntity) Run(e, de *pb.Entity) error {
	e.Meta.Locked = proto.Bool(false)
	return nil
}
