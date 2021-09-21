package interface_test

import (
	"context"
	"testing"

	"google.golang.org/protobuf/proto"

	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

func TestDropGroupCapability2(t *testing.T) {
	ctxt := context.Background()
	m, ctx := newTreeManager(t)

	dg := &pb.Group{
		Name:   proto.String("group1"),
		Number: proto.Int32(1),
		Capabilities: []pb.Capability{
			pb.Capability_CREATE_GROUP,
			pb.Capability_GLOBAL_ROOT,
		},
	}

	if err := ctx.DB.SaveGroup(ctxt, dg); err != nil {
		t.Fatal(err)
	}

	if err := m.DropGroupCapability2(ctxt, "group1", pb.Capability_CREATE_GROUP.Enum()); err != nil {
		t.Fatal(err)
	}

	g, err := ctx.DB.LoadGroup(ctxt, "group1")
	if err != nil {
		t.Fatal(err)
	}

	if len(g.GetCapabilities()) != 1 || g.GetCapabilities()[0] != pb.Capability_GLOBAL_ROOT {
		t.Error("Bad GroupCapabilities")
	}
}

func TestDropGroupCapability2BadCap(t *testing.T) {
	m, _ := newTreeManager(t)

	if err := m.DropGroupCapability2(context.Background(), "group1", nil); err != tree.ErrUnknownCapability {
		t.Error(err)
	}
}
