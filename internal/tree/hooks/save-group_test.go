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

func TestSaveGroup(t *testing.T) {
	startup.DoCallbacks()
	ctx := context.Background()

	mdb, err := db.New("memory")
	if err != nil {
		t.Fatal(err)
	}

	hook, err := NewSaveGroup(tree.WithHookStorage(mdb))
	if err != nil {
		t.Fatal(err)
	}

	g := &pb.Group{Name: proto.String("fooGroup")}

	if err := hook.Run(ctx, g, &pb.Group{}); err != nil {
		t.Fatal(err)
	}

	if _, err := mdb.LoadGroup(ctx, "fooGroup"); err != nil {
		t.Fatal("Group wasn't saved", err)
	}
}

func TestSaveGroupCB(t *testing.T) {
	saveGroupCB()
}
