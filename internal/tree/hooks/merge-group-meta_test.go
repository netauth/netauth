package hooks

import (
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

func TestMergeGroupMeta(t *testing.T) {
	hook, err := NewMergeGroupMeta(tree.RefContext{})
	if err != nil {
		t.Fatal(err)
	}

	g := &pb.Group{}
	dg := &pb.Group{
		Name:        proto.String("Unsettable Name"),
		DisplayName: proto.String("Some Group"),
	}

	if err := hook.Run(g, dg); err != nil {
		t.Fatal(err)
	}

	if g.GetName() != "" || g.GetDisplayName() != "Some Group" {
		t.Fatal("Spec error - please trace hook")
	}
}

func TestMergeGroupMetaCB(t *testing.T) {
	mergeGroupMetaCB()
}
