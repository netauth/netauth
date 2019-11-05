package hooks

import (
	"testing"

	"github.com/netauth/netauth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

func TestGroupCapabilitiesEmptyList(t *testing.T) {
	hook, err := NewSetGroupCapability(tree.RefContext{})
	if err != nil {
		t.Fatal(err)
	}

	g := &pb.Group{}
	dg := &pb.Group{}

	if err := hook.Run(g, dg); err != tree.ErrUnknownCapability {
		t.Fatal(err)
	}
}

func TestAddGroupCapabilities(t *testing.T) {
	hook, err := NewSetGroupCapability(tree.RefContext{})
	if err != nil {
		t.Fatal(err)
	}

	g := &pb.Group{}
	dg := &pb.Group{
		Capabilities: []pb.Capability{
			pb.Capability_CREATE_ENTITY,
			pb.Capability_CREATE_ENTITY,
			pb.Capability_CREATE_GROUP,
		},
	}

	if err := hook.Run(g, dg); err != nil {
		t.Fatal(err)
	}

	caps := g.GetCapabilities()
	if len(caps) != 2 || caps[0] != pb.Capability_CREATE_ENTITY {
		t.Log(g)
		t.Fatal("Spec failure, please examine hook execution")
	}
}

func TestRemoveGroupCapabilities(t *testing.T) {
	hook, err := NewRemoveGroupCapability(tree.RefContext{})
	if err != nil {
		t.Fatal(err)
	}

	g := &pb.Group{
		Capabilities: []pb.Capability{
			pb.Capability_CREATE_ENTITY,
			pb.Capability_CREATE_GROUP,
		},
	}
	dg := &pb.Group{
		Capabilities: []pb.Capability{
			pb.Capability_CREATE_ENTITY,
		},
	}

	if err := hook.Run(g, dg); err != nil {
		t.Fatal(err)
	}

	caps := g.GetCapabilities()
	if len(caps) != 1 || caps[0] != pb.Capability_CREATE_GROUP {
		t.Log(g)
		t.Fatal("Spec failure, please examine hook execution")
	}
}
