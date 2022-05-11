package interface_test

import (
	"context"
	"testing"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/protobuf/proto"

	"github.com/netauth/netauth/internal/crypto/nocrypto"
	"github.com/netauth/netauth/internal/db"
	_ "github.com/netauth/netauth/internal/db/memory"
	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"
	_ "github.com/netauth/netauth/internal/tree/hooks"

	pb "github.com/netauth/protocol"
)

func newTreeManager(t *testing.T) (*tree.Manager, tree.DB) {
	startup.DoCallbacks()

	mdb, err := db.New("memory")
	if err != nil {
		t.Fatal(err)
	}

	crypto, err := nocrypto.New(hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}

	em, err := tree.New(tree.WithStorage(mdb), tree.WithCrypto(crypto))
	if err != nil {
		t.Fatal(err)
	}

	return em, mdb
}

func addEntity(t *testing.T, db tree.DB) {
	e := &pb.Entity{
		ID:     proto.String("entity1"),
		Number: proto.Int32(1),
		Secret: proto.String("entity1"),
	}

	if err := db.SaveEntity(context.Background(), e); err != nil {
		t.Fatal(err)
	}
}

func addGroup(t *testing.T, db tree.DB) {
	g := &pb.Group{
		Name:        proto.String("group1"),
		Number:      proto.Int32(1),
		DisplayName: proto.String("Group One"),
	}

	if err := db.SaveGroup(context.Background(), g); err != nil {
		t.Fatal(err)
	}
}
