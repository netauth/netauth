package hooks

import (
	"testing"

	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

func TestAddEntityKey(t *testing.T) {
	hook, err := NewAddEntityKey(tree.RefContext{})
	if err != nil {
		t.Fatal(err)
	}

	e := &pb.Entity{Meta: &pb.EntityMeta{}}
	de := &pb.Entity{
		Meta: &pb.EntityMeta{
			Keys: []string{"KEY:code<>"},
		},
	}

	if err := hook.Run(e, de); err != nil {
		t.Fatal(err)
	}

	if e.GetMeta().GetKeys()[0] != "KEY:code<>" {
		t.Fatal("Spec error - please trace plugin")
	}
}

func TestDelEntityKey(t *testing.T) {
	hook, err := NewDelEntityKey(tree.RefContext{})
	if err != nil {
		t.Fatal(err)
	}

	e := &pb.Entity{
		Meta: &pb.EntityMeta{
			Keys: []string{"KEY:code<>", "KEY2:code<>"},
		},
	}
	de := &pb.Entity{
		Meta: &pb.EntityMeta{
			Keys: []string{"KEY:code<>"},
		},
	}

	if err := hook.Run(e, de); err != nil {
		t.Fatal(err)
	}

	if len(e.GetMeta().GetKeys()) != 1 {
		t.Log(e)
		t.Fatal("Spec error - please trace plugin")
	}
}

func TestEntityKeysCB(t *testing.T) {
	entityKeysCB()
}
