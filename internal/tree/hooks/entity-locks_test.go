package hooks

import (
	"context"
	"testing"

	"google.golang.org/protobuf/proto"

	pb "github.com/netauth/protocol"
)

func TestEntityLock(t *testing.T) {
	hook, err := NewELMLock()
	if err != nil {
		t.Fatal(err)
	}

	e := &pb.Entity{Meta: &pb.EntityMeta{}}
	de := &pb.Entity{}

	if err := hook.Run(context.Background(), e, de); err != nil {
		t.Fatal(err)
	}

	if !e.GetMeta().GetLocked() {
		t.Fatal("Spec error - please trace hook")
	}
}

func TestEntityUnLock(t *testing.T) {
	hook, err := NewELMUnlock()
	if err != nil {
		t.Fatal(err)
	}

	e := &pb.Entity{Meta: &pb.EntityMeta{Locked: proto.Bool(true)}}
	de := &pb.Entity{}

	if err := hook.Run(context.Background(), e, de); err != nil {
		t.Fatal(err)
	}

	if e.GetMeta().GetLocked() {
		t.Fatal("Spec error - please trace hook")
	}
}

func TestEntityLockCB(t *testing.T) {
	entityLockCB()
}
