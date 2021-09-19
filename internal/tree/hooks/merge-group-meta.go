package hooks

import (
	"google.golang.org/protobuf/proto"

	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
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
	startup.RegisterCallback(mergeGroupMetaCB)
}

func mergeGroupMetaCB() {
	tree.RegisterGroupHookConstructor("merge-group-meta", NewMergeGroupMeta)
}

// NewMergeGroupMeta returns a MergeGroupMeta hook configured and
// ready for use.
func NewMergeGroupMeta(c tree.RefContext) (tree.GroupHook, error) {
	return &MergeGroupMeta{tree.NewBaseHook("merge-group-meta", 50)}, nil
}
