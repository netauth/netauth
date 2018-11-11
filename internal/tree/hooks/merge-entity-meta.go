package hooks

import (
	"github.com/NetAuth/NetAuth/internal/tree"
	"github.com/golang/protobuf/proto"

	pb "github.com/NetAuth/Protocol"
)

type MergeEntityMeta struct {
	tree.BaseHook
}

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

func init() {
	tree.RegisterEntityHookConstructor("merge-entity-meta", NewMergeEntityMeta)
}

func NewMergeEntityMeta(c tree.RefContext) (tree.EntityProcessorHook, error) {
	return &MergeEntityMeta{tree.NewBaseHook("merge-entity-meta", 50)}, nil
}
