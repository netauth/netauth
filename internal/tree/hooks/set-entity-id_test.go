package hooks

import (
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/NetAuth/NetAuth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

func TestSetEntityID(t *testing.T) {
	hook, err := NewSetEntityID(tree.RefContext{})
	if err != nil {
		t.Fatal(err)
	}

	e := &pb.Entity{}
	de := &pb.Entity{ID: proto.String("entity-id")}

	if err := hook.Run(e, de); err != nil {
		t.Fatal(err)
	}
	if e.GetID() != "entity-id" {
		t.Error("Spec error - please trace hook")
	}
}
