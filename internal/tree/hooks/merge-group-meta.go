package hooks

import (
	"github.com/NetAuth/NetAuth/internal/tree"
	"github.com/golang/protobuf/proto"

	pb "github.com/NetAuth/Protocol"
)

// MergeGroupMeta provides a hook to copy metadata fields from the
// dataGroup to the group.
type MergeGroupMeta struct {
	tree.BaseHook
}

// Run attempts to copy the metadata from one group to another.
// Select fields are nil-ed out beforehand since they either require a
// specialized mechanism to edit, or a specialized capability.
func (*MergeGroupMeta) Run(g, dg *pb.Group) error {
	// There's a few fields that can't be set by merging the
	// metadata this way, so we null those out here.
	dg.Name = nil
	dg.Number = nil

	proto.Merge(g, dg)
	return nil
}

func init() {
	tree.RegisterGroupHookConstructor("merge-group-meta", NewMergeGroupMeta)
}

// NewMergeGroupMeta returns a MergeGroupMeta hook configured and
// ready for use.
func NewMergeGroupMeta(c tree.RefContext) (tree.GroupHook, error) {
	return &MergeGroupMeta{tree.NewBaseHook("merge-group-meta", 50)}, nil
}
