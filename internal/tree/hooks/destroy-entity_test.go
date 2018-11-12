package hooks

import (
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/NetAuth/NetAuth/internal/db/memdb"
	"github.com/NetAuth/NetAuth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

func TestDestroyEntity(t *testing.T) {
	mdb, err := memdb.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := tree.RefContext{
		DB: mdb,
	}

	hook, err := NewDestroyEntity(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if err = mdb.SaveEntity(&pb.Entity{ID: proto.String("foo")}); err != nil {
		t.Fatal(err)
	}
	if err = mdb.SaveEntity(&pb.Entity{ID: proto.String("bar")}); err != nil {
		t.Fatal(err)
	}

	// Act as though a delete was requested normally
	if err := hook.Run(&pb.Entity{}, &pb.Entity{ID: proto.String("foo")}); err != nil {
		t.Fatal(err)
	}

	// Act as though deleting an entity at the end of a pipeline
	if err := hook.Run(&pb.Entity{ID: proto.String("bar")}, &pb.Entity{}); err != nil {
		t.Fatal(err)
	}
}
