package hooks

import (
	"github.com/golang/protobuf/proto"

	pb "github.com/NetAuth/Protocol"
)

type LockEntity struct{}

func (*LockEntity) Name() string  { return "lock-entity" }
func (*LockEntity) Priority() int { return 50 }
func (*LockEntity) Run(e, de *pb.Entity) error {
	e.Meta.Locked = proto.Bool(true)
	return nil
}
