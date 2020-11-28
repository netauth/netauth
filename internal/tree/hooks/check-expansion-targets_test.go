package hooks

import (
	"testing"

	"github.com/netauth/netauth/internal/db"
	_ "github.com/netauth/netauth/internal/db/memory"
	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

func TestCheckExpansionTargetsDrop(t *testing.T) {
	startup.DoCallbacks()

	mdb, err := db.New("memory")
	if err != nil {
		t.Fatal(err)
	}

	hook, err := NewCheckExpansionTargets(tree.RefContext{DB: mdb})
	if err != nil {
		t.Fatal(err)
	}

	g := &pb.Group{}
	dg := &pb.Group{
		Expansions: []string{
			"DROP:deleted-group",
		},
	}

	if err := hook.Run(g, dg); err != nil {
		t.Error("Spec error - please trace hook")
	}
}

func TestCheckExpansionTargetsBad(t *testing.T) {
	startup.DoCallbacks()

	mdb, err := db.New("memory")
	if err != nil {
		t.Fatal(err)
	}

	hook, err := NewCheckExpansionTargets(tree.RefContext{DB: mdb})
	if err != nil {
		t.Fatal(err)
	}

	g := &pb.Group{}
	dg := &pb.Group{
		Expansions: []string{
			"INCLUDE:missing-group",
		},
	}

	if err := hook.Run(g, dg); err != db.ErrUnknownGroup {
		t.Error("Spec error - please trace hook")
	}
}

func TestCheckExpansionTargetsCB(t *testing.T) {
	checkExpansionTargetsCB()
}
