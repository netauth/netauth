package hooks

import (
	"testing"

	"github.com/netauth/netauth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

func TestAddEntityUM(t *testing.T) {
	hook, err := NewAddEntityUM(tree.RefContext{})
	if err != nil {
		t.Fatal(err)
	}

	e := &pb.Entity{Meta: &pb.EntityMeta{}}
	de := &pb.Entity{
		Meta: &pb.EntityMeta{
			UntypedMeta: []string{"key:value:with:colons"},
		},
	}

	if err := hook.Run(e, de); err != nil {
		t.Fatal(err)
	}

	if e.GetMeta().GetUntypedMeta()[0] != "key:value:with:colons" {
		t.Fatal("Spec error - please trace hook")
	}
}

func TestDelFuzzyEntityUM(t *testing.T) {
	hook, err := NewDelFuzzyEntityUM(tree.RefContext{})
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

	if err := hook.Run(e, de); err != nil {
		t.Fatal(err)
	}

	if len(e.GetMeta().GetUntypedMeta()) != 0 {
		t.Log(e)
		t.Fatal("Spec error - please trace hook")
	}
}

func TestDelExactEntityUM(t *testing.T) {
	hook, err := NewDelExactEntityUM(tree.RefContext{})
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

	if err := hook.Run(e, de); err != nil {
		t.Fatal(err)
	}

	if e.GetMeta().GetUntypedMeta()[0] != "key{0}:value" {
		t.Log(e)
		t.Fatal("Spec error - please trace hook")
	}
}
