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

func TestFailOnExistingGroup(t *testing.T) {
	startup.DoCallbacks()
	ctx := context.Background()

	mdb, err := db.New("memory")
	if err != nil {
		t.Fatal(err)
	}

	hook, err := NewFailOnExistingGroup(tree.WithHookStorage(mdb))
	if err != nil {
		t.Fatal(err)
	}

	err = mdb.SaveGroup(ctx, &pb.Group{
		Name: proto.String("foo"),
	})
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name    string
		wantErr error
	}{
		{"foo", tree.ErrDuplicateGroupName},
		{"bar", nil},
	}
	for i, c := range cases {
		g := &pb.Group{}
		dg := &pb.Group{Name: &c.name}
		if err := hook.Run(ctx, g, dg); err != c.wantErr {
			t.Errorf("Case %d: Got: %v Want: %v", i, err, c.wantErr)
		}
	}
}

func TestFailOnExistingGroupCB(t *testing.T) {
	failOnExistingGroupCB()
}
