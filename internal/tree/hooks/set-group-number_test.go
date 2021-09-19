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

func TestSetGroupNumber(t *testing.T) {
	startup.DoCallbacks()

	db, err := db.New("memory")
	if err != nil {
		t.Fatal(err)
	}

	hook, err := NewSetGroupNumber(tree.RefContext{DB: db})
	if err != nil {
		t.Fatal(err)
	}

	g := &pb.Group{}

	if err := hook.Run(g, &pb.Group{Number: proto.Int32(27)}); err != nil {
		t.Fatal(err)
	}

	if g.GetNumber() != 27 {
		t.Log(g)
		t.Error("Spec failure - please trace hook")
	}

	if err := hook.Run(g, &pb.Group{Number: proto.Int32(-1)}); err != nil {
		t.Fatal(err)
	}

	if g.GetNumber() != 1 {
		t.Log(g)
		t.Error("Spec failure = please trace hook")
	}
}

func TestSetGroupNumberCB(t *testing.T) {
	setGroupNumberCB()
}
