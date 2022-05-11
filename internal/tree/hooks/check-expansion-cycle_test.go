package hooks

import (
	"context"
	"testing"

	"google.golang.org/protobuf/proto"

	"github.com/netauth/netauth/internal/db"
	_ "github.com/netauth/netauth/internal/db/memory"
	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

func TestCheckExpansionCyclesDrop(t *testing.T) {
	startup.DoCallbacks()

	mdb, err := db.New("memory")
	if err != nil {
		t.Fatal(err)
	}

	hook, err := NewCheckExpansionCycles(tree.WithHookStorage(mdb))
	if err != nil {
		t.Fatal(err)
	}

	g := &pb.Group{}
	dg := &pb.Group{
		Expansions: []string{"DROP:somegroup"},
	}

	if err := hook.Run(context.Background(), g, dg); err != nil {
		t.Error(err)
	}
}

func TestCheckExpansionCycleUnknownChild(t *testing.T) {
	startup.DoCallbacks()

	mdb, err := db.New("memory")
	if err != nil {
		t.Fatal(err)
	}

	hook, err := NewCheckExpansionCycles(tree.WithHookStorage(mdb))
	if err != nil {
		t.Fatal(err)
	}

	g := &pb.Group{}
	dg := &pb.Group{
		Expansions: []string{"INCLUDE:somegroup"},
	}

	if err := hook.Run(context.Background(), g, dg); err != db.ErrUnknownGroup {
		t.Error(err)
	}
}

func TestCheckExpansionCycleCycleFound(t *testing.T) {
	startup.DoCallbacks()

	mdb, err := db.New("memory")
	if err != nil {
		t.Fatal(err)
	}

	hook, err := NewCheckExpansionCycles(tree.WithHookStorage(mdb))
	if err != nil {
		t.Fatal(err)
	}

	if err := mdb.SaveGroup(context.Background(), &pb.Group{Name: proto.String("group2"), Expansions: []string{"INCLUDE:group1"}}); err != nil {
		t.Fatal(err)
	}

	g := &pb.Group{Name: proto.String("group1")}
	dg := &pb.Group{
		Expansions: []string{"INCLUDE:group2"},
	}

	if err := hook.Run(context.Background(), g, dg); err != tree.ErrExistingExpansion {
		t.Error(err)
	}
}

func TestCheckGroupCyclesRecurser(t *testing.T) {
	startup.DoCallbacks()
	ctx := context.Background()

	mdb, err := db.New("memory")
	if err != nil {
		t.Fatal(err)
	}

	hook, err := NewCheckExpansionCycles(tree.WithHookStorage(mdb))
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
	if err := mdb.SaveGroup(ctx, grp1); err != nil {
		t.Fatal(err)
	}

	if !rhook.checkGroupCycles(ctx, grp1, "group2") {
		t.Fatal("Failed to detect direct loop")
	}

	grp2 := &pb.Group{
		Name:       proto.String("group2"),
		Expansions: []string{"INCLUDE:group3"},
	}
	if err := mdb.SaveGroup(ctx, grp2); err != nil {
		t.Fatal(err)
	}

	if !rhook.checkGroupCycles(ctx, grp1, "group4") {
		t.Fatal("Failed to error on an unloadable group")
	}

	grp3 := &pb.Group{
		Name: proto.String("group3"),
	}
	if err := mdb.SaveGroup(ctx, grp3); err != nil {
		t.Fatal(err)
	}

	if rhook.checkGroupCycles(ctx, grp1, "group4") {
		t.Fatal("Errored on an acceptable expansion")
	}
}

func TestCheckExpansionCyclesCB(t *testing.T) {
	checkExpansionCyclesCB()
}
