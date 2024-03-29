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

func TestSetEntityNumber(t *testing.T) {
	startup.DoCallbacks()
	ctx := context.Background()

	mdb, err := db.New("memory")
	if err != nil {
		t.Fatal(err)
	}

	hook, err := NewSetEntityNumber(tree.WithHookStorage(mdb))
	if err != nil {
		t.Fatal(err)
	}

	e := &pb.Entity{}
	de := &pb.Entity{Number: proto.Int32(-1)}

	if err := hook.Run(ctx, e, de); err != nil || e.GetNumber() != 1 {
		t.Log(e)
		t.Fatal(err)
	}

	de.Number = proto.Int32(27)
	if err := hook.Run(ctx, e, de); err != nil || e.GetNumber() != 27 {
		t.Log(e)
		t.Fatal(err)
	}
}

func TestSetEntityNumberCB(t *testing.T) {
	setEntityNumberCB()
}
