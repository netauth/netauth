package hooks

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-hclog"

	"github.com/netauth/netauth/internal/db"
	"github.com/netauth/netauth/internal/db/memdb"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

func TestLoadEntity(t *testing.T) {
	memdb, err := memdb.New(hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}

	hook, err := NewLoadEntity(tree.RefContext{DB: memdb})
	if err != nil {
		t.Fatal(err)
	}

	e := pb.Entity{
		ID:     proto.String("foo"),
		Number: proto.Int32(1),
		Secret: proto.String(""),
	}
	if err := memdb.SaveEntity(&e); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		ID      string
		wantErr error
	}{
		{"foo", nil},
		{"bar", db.ErrUnknownEntity},
	}

	for i, c := range cases {
		if err := hook.Run(&pb.Entity{}, &pb.Entity{ID: &c.ID}); err != c.wantErr {
			t.Errorf("Case %d: Got %v Want %v", i, err, c.wantErr)
		}
	}
}

func TestLoadEntityCB(t *testing.T) {
	loadEntityCB()
}
