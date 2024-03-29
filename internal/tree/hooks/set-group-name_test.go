package hooks

import (
	"context"
	"testing"

	"google.golang.org/protobuf/proto"

	pb "github.com/netauth/protocol"
)

func TestSetGroupName(t *testing.T) {
	hook, err := NewSetGroupName()
	if err != nil {
		t.Fatal(err)
	}

	g := &pb.Group{}
	dg := &pb.Group{Name: proto.String("fooGroup")}

	if err := hook.Run(context.Background(), g, dg); err != nil {
		t.Fatal(err)
	}

	if g.GetName() != "fooGroup" {
		t.Log(g)
		t.Error("Spec error - please trace hook")
	}
}

func TestSetGroupNameCB(t *testing.T) {
	setGroupNameCB()
}
