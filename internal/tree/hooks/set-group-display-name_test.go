package hooks

import (
	"context"
	"testing"

	"google.golang.org/protobuf/proto"

	pb "github.com/netauth/protocol"
)

func TestSetGroupDisplayName(t *testing.T) {
	hook, err := NewSetGroupDisplayName()
	if err != nil {
		t.Fatal(err)
	}

	g := &pb.Group{}
	dg := &pb.Group{DisplayName: proto.String("foo group")}

	if err := hook.Run(context.Background(), g, dg); err != nil {
		t.Fatal(err)
	}

	if g.GetDisplayName() != "foo group" {
		t.Log(g)
		t.Error("Spec failure - please trace hook")
	}
}

func TestSetGroupDisplayNameCB(t *testing.T) {
	setGroupDisplayNameCB()
}
