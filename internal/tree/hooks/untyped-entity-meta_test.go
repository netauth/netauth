package hooks

import (
	"context"
	"testing"

	pb "github.com/netauth/protocol"
)

func TestAddEntityUM(t *testing.T) {
	hook, err := NewAddEntityUM()
	if err != nil {
		t.Fatal(err)
	}

	e := &pb.Entity{Meta: &pb.EntityMeta{}}
	de := &pb.Entity{
		Meta: &pb.EntityMeta{
			UntypedMeta: []string{"key:value:with:colons"},
		},
	}

	if err := hook.Run(context.Background(), e, de); err != nil {
		t.Fatal(err)
	}

	if e.GetMeta().GetUntypedMeta()[0] != "key:value:with:colons" {
		t.Fatal("Spec error - please trace hook")
	}
}

func TestDelFuzzyEntityUM(t *testing.T) {
	hook, err := NewDelFuzzyEntityUM()
	if err != nil {
		t.Fatal(err)
	}

	e := &pb.Entity{
		Meta: &pb.EntityMeta{
			UntypedMeta: []string{
				"key{0}:value",
				"key{1}:value1",
			},
		},
	}
	de := &pb.Entity{
		Meta: &pb.EntityMeta{
			UntypedMeta: []string{"key:"},
		},
	}

	if err := hook.Run(context.Background(), e, de); err != nil {
		t.Fatal(err)
	}

	if len(e.GetMeta().GetUntypedMeta()) != 0 {
		t.Log(e)
		t.Fatal("Spec error - please trace hook")
	}
}

func TestDelExactEntityUM(t *testing.T) {
	hook, err := NewDelExactEntityUM()
	if err != nil {
		t.Fatal(err)
	}

	e := &pb.Entity{
		Meta: &pb.EntityMeta{
			UntypedMeta: []string{
				"key{0}:value",
				"key{1}:value1",
			},
		},
	}
	de := &pb.Entity{
		Meta: &pb.EntityMeta{
			UntypedMeta: []string{"key{1}:"},
		},
	}

	if err := hook.Run(context.Background(), e, de); err != nil {
		t.Fatal(err)
	}

	if e.GetMeta().GetUntypedMeta()[0] != "key{0}:value" {
		t.Log(e)
		t.Fatal("Spec error - please trace hook")
	}
}

func TestManageEntityUMCB(t *testing.T) {
	manageEntityUMCB()
}
