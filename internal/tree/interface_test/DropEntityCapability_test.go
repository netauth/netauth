package interface_test

import (
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/NetAuth/NetAuth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

func TestDropEntityCapability(t *testing.T) {
	m, ctx := newTreeManager(t)

	e := &pb.Entity{
		ID:     proto.String("entity1"),
		Number: proto.Int32(1),
		Meta: &pb.EntityMeta{
			Capabilities: []pb.Capability{
				pb.Capability_CREATE_ENTITY,
				pb.Capability_GLOBAL_ROOT,
			},
		},
	}
	if err := ctx.DB.SaveEntity(e); err != nil {
		t.Fatal(err)
	}

	if err := m.DropEntityCapability("entity1", "GLOBAL_ROOT"); err != nil {
		t.Error(err)
	}

	e, err := ctx.DB.LoadEntity("entity1")
	if err != nil {
		t.Fatal(err)
	}

	caps := e.GetMeta().GetCapabilities()
	if len(caps) != 1 || caps[0] != pb.Capability_CREATE_ENTITY {
		t.Error("Capability not correctly removed")
	}
}

func TestDropEntityCapabilityUnknownCapability(t *testing.T) {
	m, _ := newTreeManager(t)

	if err := m.DropEntityCapability("entity1", "UNKNOWN"); err != tree.ErrUnknownCapability {
		t.Error(err)
	}
}
