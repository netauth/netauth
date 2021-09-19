package hooks

import (
	"testing"

	"google.golang.org/protobuf/proto"

	"github.com/netauth/netauth/internal/db"
	_ "github.com/netauth/netauth/internal/db/memory"
	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

func TestDestroyEntity(t *testing.T) {
	startup.DoCallbacks()

	mdb, err := db.New("memory")
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

func TestDestroyEntityCB(t *testing.T) {
	destroyEntityCB()
}
