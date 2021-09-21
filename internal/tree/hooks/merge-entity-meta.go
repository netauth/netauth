package hooks

import (
	"context"

	"google.golang.org/protobuf/proto"

	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

// MergeEntityMeta provides a hook to copy metadata fields from the
// dataEntity to the entity.
type MergeEntityMeta struct {
	tree.BaseHook
}

// Run attempts to copy the metadata from one entity to another.
// Select fields are nil-ed out beforehand since they either require a
// specialized mechanism to edit, or a specialized capability.
func (*MergeEntityMeta) Run(_ context.Context, e, de *pb.Entity) error {
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
	startup.RegisterCallback(mergeEntityMetaCB)
}

func mergeEntityMetaCB() {
	tree.RegisterEntityHookConstructor("merge-entity-meta", NewMergeEntityMeta)
}

// NewMergeEntityMeta returns a MergeEntityMeta hook configured and
// ready for use.
func NewMergeEntityMeta(c tree.RefContext) (tree.EntityHook, error) {
	return &MergeEntityMeta{tree.NewBaseHook("merge-entity-meta", 50)}, nil
}
