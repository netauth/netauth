package hooks

import (
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/netauth/netauth/internal/db/memdb"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

func TestSetGroupNumber(t *testing.T) {
	db, err := memdb.New()
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
