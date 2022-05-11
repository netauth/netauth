package hooks

import (
	"context"
	"testing"

	pb "github.com/netauth/protocol"
)

func TestPatchGroupExpansionsInclude(t *testing.T) {
	hook, err := NewPatchGroupExpansions()
	if err != nil {
		t.Fatal(err)
	}

	g := &pb.Group{}
	dg := &pb.Group{
		Expansions: []string{
			"INCLUDE:group1",
		},
	}

	if err := hook.Run(context.Background(), g, dg); err != nil {
		t.Fatal(err)
	}

	if len(g.GetExpansions()) != 1 || g.GetExpansions()[0] != "INCLUDE:group1" {
		t.Error("Spec error - please trace hook")
	}
}

func TestPatchGroupExpansionsDrop(t *testing.T) {
	hook, err := NewPatchGroupExpansions()
	if err != nil {
		t.Fatal(err)
	}

	g := &pb.Group{
		Expansions: []string{
			"INCLUDE:group1",
			"EXCLUDE:group2",
		},
	}
	dg := &pb.Group{
		Expansions: []string{
			"DROP:group1",
		},
	}

	if err := hook.Run(context.Background(), g, dg); err != nil {
		t.Fatal(err)
	}

	if len(g.GetExpansions()) != 1 || g.GetExpansions()[0] != "EXCLUDE:group2" {
		t.Error("Spec error - please trace hook")
	}

}

func TestPatchGroupExpansionsCB(t *testing.T) {
	patchGroupExpansionsCB()
}
