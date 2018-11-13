package hooks

import (
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/NetAuth/NetAuth/internal/db/memdb"
	"github.com/NetAuth/NetAuth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

func TestSetEntityNumber(t *testing.T) {
	memdb, err := memdb.New()
	if err != nil {
		t.Fatal(err)
	}

	hook, err := NewSetEntityNumber(tree.RefContext{DB: memdb})
	if err != nil {
		t.Fatal(err)
	}

	e := &pb.Entity{}
	de := &pb.Entity{Number: proto.Int32(-1)}

	if err := hook.Run(e, de); err != nil || e.GetNumber() != 1 {
		t.Log(e)
		t.Fatal(err)
	}

	de.Number = proto.Int32(27)
	if err := hook.Run(e, de); err != nil || e.GetNumber() != 27 {
		t.Log(e)
		t.Fatal(err)
	}
}
