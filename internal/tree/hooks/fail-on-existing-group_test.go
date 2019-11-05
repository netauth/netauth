package hooks

import (
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/netauth/netauth/internal/db/memdb"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

func TestFailOnExistingGroup(t *testing.T) {
	mdb, err := memdb.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := tree.RefContext{
		DB: mdb,
	}

	hook, err := NewFailOnExistingGroup(ctx)
	if err != nil {
		t.Fatal(err)
	}

	err = mdb.SaveGroup(&pb.Group{
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
		if err := hook.Run(g, dg); err != c.wantErr {
			t.Errorf("Case %d: Got: %v Want: %v", i, err, c.wantErr)
		}
	}
}
