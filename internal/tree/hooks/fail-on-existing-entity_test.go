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

func TestFailOnExistingEntity(t *testing.T) {
	startup.DoCallbacks()
	ctx := context.Background()

	mdb, err := db.New("memory")
	if err != nil {
		t.Fatal(err)
	}
	rctx := tree.RefContext{
		DB: mdb,
	}

	hook, err := NewFailOnExistingEntity(rctx)
	if err != nil {
		t.Fatal(err)
	}

	err = mdb.SaveEntity(ctx, &pb.Entity{
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
		if err := hook.Run(ctx, e, de); err != c.wantErr {
			t.Errorf("Case %d: Got: %v Want: %v", i, err, c.wantErr)
		}
	}
}

func TestFailOnExistingEntityCB(t *testing.T) {
	failOnExistingEntityCB()
}
