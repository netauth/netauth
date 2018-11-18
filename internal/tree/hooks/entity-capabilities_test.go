package hooks

import (
	"testing"

	"github.com/NetAuth/NetAuth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

func TestAddEntityCapabilities(t *testing.T) {
	hook, err := NewSetEntityCapability(tree.RefContext{})
	if err != nil {
		t.Fatal(err)
	}

	e := &pb.Entity{Meta: &pb.EntityMeta{}}
	de := &pb.Entity{
		Meta: &pb.EntityMeta{
			Capabilities: []pb.Capability{
				pb.Capability_CREATE_ENTITY,
				pb.Capability_CREATE_ENTITY,
				pb.Capability_CREATE_GROUP,
			},
		},
	}

	if err := hook.Run(e, de); err != nil {
		t.Fatal(err)
	}

	caps := e.GetMeta().GetCapabilities()
	if len(caps) != 2 || caps[0] != pb.Capability_CREATE_ENTITY {
		t.Log(e)
		t.Fatal("Spec failure, please examine hook execution")
	}
}

func TestRemoveEntityCapabilities(t *testing.T) {
	hook, err := NewRemoveEntityCapability(tree.RefContext{})
	if err != nil {
		t.Fatal(err)
	}

	e := &pb.Entity{
		Meta: &pb.EntityMeta{
			Capabilities: []pb.Capability{
				pb.Capability_CREATE_ENTITY,
				pb.Capability_CREATE_GROUP,
			},
		},
	}
	de := &pb.Entity{
		Meta: &pb.EntityMeta{
			Capabilities: []pb.Capability{
				pb.Capability_CREATE_ENTITY,
			},
		},
	}

	if err := hook.Run(e, de); err != nil {
		t.Fatal(err)
	}

	caps := e.GetMeta().GetCapabilities()
	if len(caps) != 1 || caps[0] != pb.Capability_CREATE_GROUP {
		t.Log(e)
		t.Fatal("Spec failure, please examine hook execution")
	}
}
