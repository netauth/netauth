package hooks

import (
	"context"
	"testing"

	"google.golang.org/protobuf/proto"

	pb "github.com/netauth/protocol"
)

func TestMergeGroupMeta(t *testing.T) {
	hook, err := NewMergeGroupMeta()
	if err != nil {
		t.Fatal(err)
	}

	g := &pb.Group{}
	dg := &pb.Group{
		Name:        proto.String("Unsettable Name"),
		DisplayName: proto.String("Some Group"),
	}

	if err := hook.Run(context.Background(), g, dg); err != nil {
		t.Fatal(err)
	}

	if g.GetName() != "" || g.GetDisplayName() != "Some Group" {
		t.Fatal("Spec error - please trace hook")
	}
}

func TestMergeGroupMetaCB(t *testing.T) {
	mergeGroupMetaCB()
}
