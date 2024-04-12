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

func TestFailOnEmptyEntity(t *testing.T) {
	startup.DoCallbacks()
	ctx := context.Background()

	mdb, err := db.New("memory")
	if err != nil {
		t.Fatal(err)
	}

	hook, err := NewFailOnEmptyEntity(tree.WithHookStorage(mdb))
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
		{"foo", nil},
		{"bar", nil},
		{"", db.ErrNoValue},
	}
	for i, c := range cases {
		e := &pb.Entity{}
		de := &pb.Entity{ID: &c.ID}
		if err := hook.Run(ctx, e, de); err != c.wantErr {
			t.Errorf("Case %d: Got: %v Want: %v", i, err, c.wantErr)
		}
	}
}

func TestFailOnEmptyEntityCB(t *testing.T) {
	failOnEmptyEntityCB()
}
