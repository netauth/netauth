package hooks

import (
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/netauth/netauth/internal/db/memdb"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

func TestDestroyGroup(t *testing.T) {
	mdb, err := memdb.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := tree.RefContext{
		DB: mdb,
	}

	hook, err := NewDestroyGroup(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if err = mdb.SaveGroup(&pb.Group{Name: proto.String("foo")}); err != nil {
		t.Fatal(err)
	}
	if err = mdb.SaveGroup(&pb.Group{Name: proto.String("bar")}); err != nil {
		t.Fatal(err)
	}

	// Act as though a delete was requested normally
	if err := hook.Run(&pb.Group{}, &pb.Group{Name: proto.String("foo")}); err != nil {
		t.Fatal(err)
	}

	// Act as though deleting an entity at the end of a pipeline
	if err := hook.Run(&pb.Group{Name: proto.String("bar")}, &pb.Group{}); err != nil {
		t.Fatal(err)
	}
}
