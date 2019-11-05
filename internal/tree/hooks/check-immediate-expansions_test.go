package hooks

import (
	"testing"

	"github.com/netauth/netauth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

func TestCheckImmediateExpansionsDrop(t *testing.T) {
	hook, err := NewCheckImmediateExpansions(tree.RefContext{})
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

	if err := hook.Run(g, dg); err != nil {
		t.Error("Spec error - please trace hook")
	}
}

func TestCheckImmediateExpansionsExisting(t *testing.T) {
	hook, err := NewCheckImmediateExpansions(tree.RefContext{})
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

	if err := hook.Run(g, dg); err != tree.ErrExistingExpansion {
		t.Error("Spec error - please trace hook")
	}
}
