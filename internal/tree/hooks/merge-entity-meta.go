package hooks

import (
	"github.com/golang/protobuf/proto"

	pb "github.com/NetAuth/Protocol"
)

type MergeEntityMeta struct{}

func (*MergeEntityMeta) Name() string  { return "merge-entity-meta" }
func (*MergeEntityMeta) Priority() int { return 50 }
func (*MergeEntityMeta) Run(e, de *pb.Entity) error {
	// There's a few fields that can't be set by merging the
	// metadata this way, so we null those out here.
	de.Meta.Capabilities = nil
	de.Meta.Groups = nil
	de.Meta.Keys = nil
	de.Meta.UntypedMeta = nil

	proto.Merge(e, de)
	return nil
}
