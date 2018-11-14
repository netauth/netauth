package hooks

import (
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/NetAuth/NetAuth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

func TestSetGroupDisplayName(t *testing.T) {
	hook, err := NewSetGroupDisplayName(tree.RefContext{})
	if err != nil {
		t.Fatal(err)
	}

	g := &pb.Group{}
	dg := &pb.Group{DisplayName: proto.String("foo group")}

	if err := hook.Run(g, dg); err != nil {
		t.Fatal(err)
	}

	if g.GetDisplayName() != "foo group" {
		t.Log(g)
		t.Error("Spec failure - please trace hook")
	}
}
