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

func TestDestroyGroup(t *testing.T) {
	startup.DoCallbacks()
	ctx := context.Background()

	mdb, err := db.New("memory")
	if err != nil {
		t.Fatal(err)
	}

	hook, err := NewDestroyGroup(tree.WithHookStorage(mdb))
	if err != nil {
		t.Fatal(err)
	}

	if err = mdb.SaveGroup(ctx, &pb.Group{Name: proto.String("foo")}); err != nil {
		t.Fatal(err)
	}
	if err = mdb.SaveGroup(ctx, &pb.Group{Name: proto.String("bar")}); err != nil {
		t.Fatal(err)
	}

	// Act as though a delete was requested normally
	if err := hook.Run(ctx, &pb.Group{}, &pb.Group{Name: proto.String("foo")}); err != nil {
		t.Fatal(err)
	}

	// Act as though deleting an entity at the end of a pipeline
	if err := hook.Run(ctx, &pb.Group{Name: proto.String("bar")}, &pb.Group{}); err != nil {
		t.Fatal(err)
	}
}

func TestDestroyGroupCB(t *testing.T) {
	destroyGroupCB()
}
