package hooks

import (
	"context"
	"testing"

	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

func TestAddGroupUM(t *testing.T) {
	hook, err := NewAddGroupUM(tree.RefContext{})
	if err != nil {
		t.Fatal(err)
	}

	e := &pb.Group{}
	de := &pb.Group{
		UntypedMeta: []string{"key:value:with:colons"},
	}

	if err := hook.Run(context.Background(), e, de); err != nil {
		t.Fatal(err)
	}

	if e.GetUntypedMeta()[0] != "key:value:with:colons" {
		t.Fatal("Spec error - please trace hook")
	}
}

func TestDelFuzzyGroupUM(t *testing.T) {
	hook, err := NewDelFuzzyGroupUM(tree.RefContext{})
	if err != nil {
		t.Fatal(err)
	}

	e := &pb.Group{
		UntypedMeta: []string{
			"key{0}:value",
			"key{1}:value1",
		},
	}
	de := &pb.Group{
		UntypedMeta: []string{"key:"},
	}

	if err := hook.Run(context.Background(), e, de); err != nil {
		t.Fatal(err)
	}

	if len(e.GetUntypedMeta()) != 0 {
		t.Log(e)
		t.Fatal("Spec error - please trace hook")
	}
}

func TestDelExactGroupUM(t *testing.T) {
	hook, err := NewDelExactGroupUM(tree.RefContext{})
	if err != nil {
		t.Fatal(err)
	}

	e := &pb.Group{
		UntypedMeta: []string{
			"key{0}:value",
			"key{1}:value1",
		},
	}
	de := &pb.Group{
		UntypedMeta: []string{"key{1}:"},
	}

	if err := hook.Run(context.Background(), e, de); err != nil {
		t.Fatal(err)
	}

	if e.GetUntypedMeta()[0] != "key{0}:value" {
		t.Log(e)
		t.Fatal("Spec error - please trace hook")
	}
}

func TestManageGroupUMCB(t *testing.T) {
	manageGroupUMCB()
}
