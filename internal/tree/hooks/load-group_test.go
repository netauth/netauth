package hooks

import (
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/netauth/netauth/internal/db"
	_ "github.com/netauth/netauth/internal/db/memory"
	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

func TestLoadGroup(t *testing.T) {
	startup.DoCallbacks()

	memdb, err := db.New("memory")
	if err != nil {
		t.Fatal(err)
	}

	hook, err := NewLoadGroup(tree.RefContext{DB: memdb})
	if err != nil {
		t.Fatal(err)
	}

	g := &pb.Group{
		Name:   proto.String("group"),
		Number: proto.Int32(1),
	}
	if err := memdb.SaveGroup(g); err != nil {
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
		if err := hook.Run(&pb.Group{}, &pb.Group{Name: &c.Name}); err != c.wantErr {
			t.Errorf("Case %d: Got %v Want %v", i, err, c.wantErr)
		}
	}
}

func TestLoadGroupCB(t *testing.T) {
	loadGroupCB()
}
