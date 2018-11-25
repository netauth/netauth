package interface_test

import (
	"testing"

	"github.com/NetAuth/NetAuth/internal/tree"
	
	pb "github.com/NetAuth/Protocol"
)

func TestSetGroupCapability(t *testing.T) {
	m, ctx := newTreeManager(t)

	addGroup(t, ctx)

	if err := m.SetGroupCapabilityByName("group1", "GLOBAL_ROOT"); err != nil {
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

func TestSetGroupCapabilityBadCap(t *testing.T) {
	m, _ := newTreeManager(t)

	if err := m.SetGroupCapabilityByName("group1", "UNKNOWN"); err != tree.ErrUnknownCapability {
		t.Error(err)
	}
}
