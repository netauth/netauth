package interface_test

import (
	"context"
	"testing"

	"google.golang.org/protobuf/proto"

	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

func TestDropEntityCapability2(t *testing.T) {
	ctxt := context.Background()
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
	if err := ctx.DB.SaveEntity(ctxt, e); err != nil {
		t.Fatal(err)
	}

	if err := m.DropEntityCapability2(ctxt, "entity1", pb.Capability_GLOBAL_ROOT.Enum()); err != nil {
		t.Error(err)
	}

	e, err := ctx.DB.LoadEntity(ctxt, "entity1")
	if err != nil {
		t.Fatal(err)
	}

	caps := e.GetMeta().GetCapabilities()
	if len(caps) != 1 || caps[0] != pb.Capability_CREATE_ENTITY {
		t.Error("Capability not correctly removed")
	}
}

func TestDropEntityCapability2UnknownCapability(t *testing.T) {
	m, _ := newTreeManager(t)

	if err := m.DropEntityCapability2(context.Background(), "entity1", nil); err != tree.ErrUnknownCapability {
		t.Error(err)
	}
}
