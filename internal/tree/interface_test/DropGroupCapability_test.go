package interface_test

import (
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/NetAuth/NetAuth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

func TestDropGroupCapability(t *testing.T) {
	m, ctx := newTreeManager(t)

	dg := &pb.Group{
		Name:   proto.String("group1"),
		Number: proto.Int32(1),
		Capabilities: []pb.Capability{
			pb.Capability_CREATE_GROUP,
			pb.Capability_GLOBAL_ROOT,
		},
	}

	if err := ctx.DB.SaveGroup(dg); err != nil {
		t.Fatal(err)
	}

	if err := m.DropGroupCapability("group1", "CREATE_GROUP"); err != nil {
		t.Fatal(err)
	}

	g, err := ctx.DB.LoadGroup("group1")
	if err != nil {
		t.Fatal(err)
	}

	if len(g.GetCapabilities()) != 1 || g.GetCapabilities()[0] != pb.Capability_GLOBAL_ROOT {
		t.Error("Bad GroupCapabilities")
	}
}

func TestDropGroupCapabilityBadCap(t *testing.T) {
	m, _ := newTreeManager(t)

	if err := m.DropGroupCapability("group1", "UNKNOWN"); err != tree.ErrUnknownCapability {
		t.Error(err)
	}
}
