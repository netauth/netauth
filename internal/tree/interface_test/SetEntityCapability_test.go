package interface_test

import (
	"context"
	"testing"

	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

func TestSetEntityCapability2(t *testing.T) {
	ctxt := context.Background()
	m, ctx := newTreeManager(t)

	addEntity(t, ctx)

	if err := m.SetEntityCapability2(ctxt, "entity1", pb.Capability_GLOBAL_ROOT.Enum()); err != nil {
		t.Error(err)
	}

	e, err := ctx.DB.LoadEntity(ctxt, "entity1")
	if err != nil {
		t.Fatal(err)
	}

	if e.GetMeta().GetCapabilities()[0] != pb.Capability_GLOBAL_ROOT {
		t.Error("Capability not assigned")
	}
}

func TestSetEntityCapability2UnknownCapability(t *testing.T) {
	m, _ := newTreeManager(t)

	if err := m.SetEntityCapability2(context.Background(), "entity1", nil); err != tree.ErrUnknownCapability {
		t.Error(err)
	}
}
