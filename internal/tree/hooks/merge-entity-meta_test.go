package hooks

import (
	"context"
	"testing"

	"google.golang.org/protobuf/proto"

	pb "github.com/netauth/protocol"
)

func TestMergeEntityMeta(t *testing.T) {
	hook, err := NewMergeEntityMeta()
	if err != nil {
		t.Fatal(err)
	}

	e := &pb.Entity{Meta: &pb.EntityMeta{}}
	de := &pb.Entity{
		Meta: &pb.EntityMeta{
			GECOS:  proto.String("PFY"),
			Groups: []string{"not-to-be-merged"},
		},
	}

	if err := hook.Run(context.Background(), e, de); err != nil {
		t.Fatal(err)
	}

	if len(e.GetMeta().GetGroups()) != 0 || e.GetMeta().GetGECOS() != "PFY" {
		t.Fatal("Spec error - please trace hook")
	}
}

func TestMergeEntityMetaCB(t *testing.T) {
	mergeEntityMetaCB()
}
