package hooks

import (
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/netauth/netauth/internal/db/memdb"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

func TestSaveEntity(t *testing.T) {
	mdb, err := memdb.New()
	if err != nil {
		t.Fatal(err)
	}

	hook, err := NewSaveEntity(tree.RefContext{DB: mdb})
	if err != nil {
		t.Fatal(err)
	}

	e := &pb.Entity{ID: proto.String("foobar")}

	if err := hook.Run(e, &pb.Entity{}); err != nil {
		t.Fatal(err)
	}

	if _, err := mdb.LoadEntity("foobar"); err != nil {
		t.Fatal("Entity wasn't saved", err)
	}
}
