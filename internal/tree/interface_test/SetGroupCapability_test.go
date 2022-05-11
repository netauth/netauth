package interface_test

import (
	"context"
	"testing"

	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

func TestSetGroupCapability2(t *testing.T) {
	ctxt := context.Background()
	m, mdb := newTreeManager(t)

	addGroup(t, mdb)

	if err := m.SetGroupCapability2(ctxt, "group1", pb.Capability_GLOBAL_ROOT.Enum()); err != nil {
		t.Fatal(err)
	}

	g, err := mdb.LoadGroup(ctxt, "group1")
	if err != nil {
		t.Fatal(err)
	}

	if len(g.GetCapabilities()) != 1 || g.GetCapabilities()[0] != pb.Capability_GLOBAL_ROOT {
		t.Error("Bad GroupCapabilities")
	}
}

func TestSetGroupCapability2BadCap(t *testing.T) {
	m, _ := newTreeManager(t)

	if err := m.SetGroupCapability2(context.Background(), "group1", nil); err != tree.ErrUnknownCapability {
		t.Error(err)
	}
}
