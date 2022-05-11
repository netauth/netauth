package hooks

import (
	"context"
	"testing"

	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

func TestCheckImmediateExpansionsDrop(t *testing.T) {
	hook, err := NewCheckImmediateExpansions()
	if err != nil {
		t.Fatal(err)
	}

	g := &pb.Group{
		Expansions: []string{
			"INCLUDE:irrelevant",
		},
	}
	dg := &pb.Group{
		Expansions: []string{
			"DROP:foo1",
			"DROP:foo2",
		},
	}

	if err := hook.Run(context.Background(), g, dg); err != nil {
		t.Error("Spec error - please trace hook")
	}
}

func TestCheckImmediateExpansionsExisting(t *testing.T) {
	hook, err := NewCheckImmediateExpansions()
	if err != nil {
		t.Fatal(err)
	}

	g := &pb.Group{
		Expansions: []string{
			"INCLUDE:foo2",
		},
	}
	dg := &pb.Group{
		Expansions: []string{
			"INCLUDE:foo1",
			"INCLUDE:foo2",
		},
	}

	if err := hook.Run(context.Background(), g, dg); err != tree.ErrExistingExpansion {
		t.Error("Spec error - please trace hook")
	}
}

func TestCheckImmediateExpansionsCB(t *testing.T) {
	checkImmediateExpansionsCB()
}
