package interface_test

import (
	"testing"

	"github.com/NetAuth/NetAuth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

func TestSetEntityCapability(t *testing.T) {
	m, ctx := newTreeManager(t)

	addEntity(t, ctx)

	if err := m.SetEntityCapability("entity1", "GLOBAL_ROOT"); err != nil {
		t.Error(err)
	}

	e, err := ctx.DB.LoadEntity("entity1")
	if err != nil {
		t.Fatal(err)
	}

	if e.GetMeta().GetCapabilities()[0] != pb.Capability_GLOBAL_ROOT {
		t.Error("Capability not assigned")
	}
}

func TestSetEntityCapabilityUnknownCapability(t *testing.T) {
	m, _ := newTreeManager(t)

	if err := m.SetEntityCapability("entity1", "UNKNOWN"); err != tree.ErrUnknownCapability {
		t.Error(err)
	}
}
