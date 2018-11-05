package hooks

import (
	pb "github.com/NetAuth/Protocol"
)

type EnsureEntityMeta struct{}

func (*EnsureEntityMeta) Name() string  { return "ensure-entity-meta" }
func (*EnsureEntityMeta) Priority() int { return 25 }
func (*EnsureEntityMeta) Run(e, de *pb.Entity) error {
	if e.Meta == nil {
		e.Meta = &pb.EntityMeta{}
	}
	return nil
}
