package hooks

import (
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/NetAuth/NetAuth/internal/db"
	"github.com/NetAuth/NetAuth/internal/db/memdb"
	"github.com/NetAuth/NetAuth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

func TestLoadEntity(t *testing.T) {
	memdb, err := memdb.New()
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
