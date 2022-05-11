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

func TestSetManagingGroup(t *testing.T) {
	startup.DoCallbacks()
	ctx := context.Background()

	mdb, err := db.New("memory")
	if err != nil {
		t.Fatal(err)
	}

	if err := mdb.SaveGroup(ctx, &pb.Group{Name: proto.String("bar")}); err != nil {
		t.Fatal(err)
	}

	hook, err := NewSetManagingGroup(tree.WithHookStorage(mdb))
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name          string
		managedby     string
		wantErr       error
		wantManagedBy string
	}{
		{"foo", "", nil, ""},
		{"foo", "foo", nil, "foo"},
		{"foo", "baz", db.ErrUnknownGroup, ""},
		{"foo", "bar", nil, "bar"},
	}

	for i, c := range cases {
		g := &pb.Group{}
		dg := &pb.Group{
			Name:      proto.String(c.name),
			ManagedBy: proto.String(c.managedby),
		}
		if err := hook.Run(ctx, g, dg); err != c.wantErr {
			t.Errorf("Case %d: Got %v Want %v", i, err, c.wantErr)
		}
		if g.GetManagedBy() != c.wantManagedBy {
			t.Errorf("Case %d: spec error - please trace hook", i)
		}
	}
}

func TestSetManagingGroupCB(t *testing.T) {
	setManagingGroupCB()
}
