package hooks

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-hclog"

	"github.com/netauth/netauth/internal/db"
	"github.com/netauth/netauth/internal/db/memdb"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

func TestCheckExpansionCyclesDrop(t *testing.T) {
	memdb, err := memdb.New(hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}

	hook, err := NewCheckExpansionCycles(tree.RefContext{DB: memdb})
	if err != nil {
		t.Fatal(err)
	}

	g := &pb.Group{}
	dg := &pb.Group{
		Expansions: []string{"DROP:somegroup"},
	}

	if err := hook.Run(g, dg); err != nil {
		t.Error(err)
	}
}

func TestCheckExpansionCycleUnknownChild(t *testing.T) {
	memdb, err := memdb.New(hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}

	hook, err := NewCheckExpansionCycles(tree.RefContext{DB: memdb})
	if err != nil {
		t.Fatal(err)
	}

	g := &pb.Group{}
	dg := &pb.Group{
		Expansions: []string{"INCLUDE:somegroup"},
	}

	if err := hook.Run(g, dg); err != db.ErrUnknownGroup {
		t.Error(err)
	}
}

func TestCheckExpansionCycleCycleFound(t *testing.T) {
	memdb, err := memdb.New(hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}

	hook, err := NewCheckExpansionCycles(tree.RefContext{DB: memdb})
	if err != nil {
		t.Fatal(err)
	}

	if err := memdb.SaveGroup(&pb.Group{Name: proto.String("group2"), Expansions: []string{"INCLUDE:group1"}}); err != nil {
		t.Fatal(err)
	}

	g := &pb.Group{Name: proto.String("group1")}
	dg := &pb.Group{
		Expansions: []string{"INCLUDE:group2"},
	}

	if err := hook.Run(g, dg); err != tree.ErrExistingExpansion {
		t.Error(err)
	}
}

func TestCheckGroupCyclesRecurser(t *testing.T) {
	memdb, err := memdb.New(hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}

	hook, err := NewCheckExpansionCycles(tree.RefContext{DB: memdb})
	if err != nil {
		t.Fatal(err)
	}

	rhook, ok := hook.(*CheckExpansionCycles)
	if !ok {
		t.Fatal("type error")
	}

	grp1 := &pb.Group{
		Name:       proto.String("group1"),
		Expansions: []string{"INCLUDE:group2"},
	}
	if err := memdb.SaveGroup(grp1); err != nil {
		t.Fatal(err)
	}

	if !rhook.checkGroupCycles(grp1, "group2") {
		t.Fatal("Failed to detect direct loop")
	}

	grp2 := &pb.Group{
		Name:       proto.String("group2"),
		Expansions: []string{"INCLUDE:group3"},
	}
	if err := memdb.SaveGroup(grp2); err != nil {
		t.Fatal(err)
	}

	if !rhook.checkGroupCycles(grp1, "group4") {
		t.Fatal("Failed to error on an unloadable group")
	}

	grp3 := &pb.Group{
		Name: proto.String("group3"),
	}
	if err := memdb.SaveGroup(grp3); err != nil {
		t.Fatal(err)
	}

	if rhook.checkGroupCycles(grp1, "group4") {
		t.Fatal("Errored on an acceptable expansion")
	}
}
