package hooks

import (
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

func TestEntityLock(t *testing.T) {
	hook, err := NewELMLock(tree.RefContext{})
	if err != nil {
		t.Fatal(err)
	}

	e := &pb.Entity{Meta: &pb.EntityMeta{}}
	de := &pb.Entity{}

	if err := hook.Run(e, de); err != nil {
		t.Fatal(err)
	}

	if !e.GetMeta().GetLocked() {
		t.Fatal("Spec error - please trace hook")
	}
}

func TestEntityUnLock(t *testing.T) {
	hook, err := NewELMUnlock(tree.RefContext{})
	if err != nil {
		t.Fatal(err)
	}

	e := &pb.Entity{Meta: &pb.EntityMeta{Locked: proto.Bool(true)}}
	de := &pb.Entity{}

	if err := hook.Run(e, de); err != nil {
		t.Fatal(err)
	}

	if e.GetMeta().GetLocked() {
		t.Fatal("Spec error - please trace hook")
	}
}
