package hooks

import (
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/NetAuth/NetAuth/internal/db/memdb"
	"github.com/NetAuth/NetAuth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

func TestFailOnExistingEntity(t *testing.T) {
	mdb, err := memdb.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := tree.RefContext{
		DB: mdb,
	}

	hook, err := NewFailOnExistingEntity(ctx)
	if err != nil {
		t.Fatal(err)
	}

	err = mdb.SaveEntity(&pb.Entity{
		ID:     proto.String("foo"),
		Number: proto.Int32(42),
	})
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		ID      string
		wantErr error
	}{
		{"foo", tree.ErrDuplicateEntityID},
		{"bar", nil},
	}
	for i, c := range cases {
		e := &pb.Entity{}
		de := &pb.Entity{ID: &c.ID}
		if err := hook.Run(e, de); err != c.wantErr {
			t.Errorf("Case %d: Got: %v Want: %v", i, err, c.wantErr)
		}
	}
}
