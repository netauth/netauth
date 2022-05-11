package hooks

import (
	"context"
	"testing"

	"google.golang.org/protobuf/proto"

	pb "github.com/netauth/protocol"
)

func TestSetEntityID(t *testing.T) {
	hook, err := NewSetEntityID()
	if err != nil {
		t.Fatal(err)
	}

	e := &pb.Entity{}
	de := &pb.Entity{ID: proto.String("entity-id")}

	if err := hook.Run(context.Background(), e, de); err != nil {
		t.Fatal(err)
	}
	if e.GetID() != "entity-id" {
		t.Error("Spec error - please trace hook")
	}
}

func TestSetEntityIDCB(t *testing.T) {
	setEntityIDCB()
}
