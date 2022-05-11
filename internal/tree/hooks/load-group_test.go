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

func TestLoadGroup(t *testing.T) {
	startup.DoCallbacks()
	ctx := context.Background()

	mdb, err := db.New("memory")
	if err != nil {
		t.Fatal(err)
	}

	hook, err := NewLoadGroup(tree.WithHookStorage(mdb))
	if err != nil {
		t.Fatal(err)
	}

	g := &pb.Group{
		Name:   proto.String("group"),
		Number: proto.Int32(1),
	}
	if err := mdb.SaveGroup(ctx, g); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		Name    string
		wantErr error
	}{
		{"group", nil},
		{"unknown", db.ErrUnknownGroup},
	}

	for i, c := range cases {
		if err := hook.Run(ctx, &pb.Group{}, &pb.Group{Name: &c.Name}); err != c.wantErr {
			t.Errorf("Case %d: Got %v Want %v", i, err, c.wantErr)
		}
	}
}

func TestLoadGroupCB(t *testing.T) {
	loadGroupCB()
}
