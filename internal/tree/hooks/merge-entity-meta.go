package hooks

import (
	"github.com/netauth/netauth/internal/tree"
	"github.com/golang/protobuf/proto"

	pb "github.com/NetAuth/Protocol"
)

// MergeEntityMeta provides a hook to copy metadata fields from the
// dataEntity to the entity.
type MergeEntityMeta struct {
	tree.BaseHook
}

// Run attempts to copy the metadata from one entity to another.
// Select fields are nil-ed out beforehand since they either require a
// specialized mechanism to edit, or a specialized capability.
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

// NewMergeEntityMeta returns a MergeEntityMeta hook configured and
// ready for use.
func NewMergeEntityMeta(c tree.RefContext) (tree.EntityHook, error) {
	return &MergeEntityMeta{tree.NewBaseHook("merge-entity-meta", 50)}, nil
}
