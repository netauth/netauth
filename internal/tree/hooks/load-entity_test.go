package hooks

import (
	"testing"

	"google.golang.org/protobuf/proto"

	"github.com/netauth/netauth/internal/db"
	_ "github.com/netauth/netauth/internal/db/memory"
	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"

	pb "github.com/netauth/protocol"
)

func TestLoadEntity(t *testing.T) {
	startup.DoCallbacks()

	memdb, err := db.New("memory")
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
