package interface_test

import (
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

func newTreeManager(t *testing.T) (*tree.Manager, tree.RefContext) {
	startup.DoCallbacks()

	mdb, err := db.New("memory")
	if err != nil {
		t.Fatal(err)
	}

	crypto, err := nocrypto.New(hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}

	ctx := tree.RefContext{
		DB:     mdb,
		Crypto: crypto,
	}

	em, err := tree.New(ctx.DB, ctx.Crypto, hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}

	return em, ctx
}

func addEntity(t *testing.T, ctx tree.RefContext) {
	e := &pb.Entity{
		ID:     proto.String("entity1"),
		Number: proto.Int32(1),
		Secret: proto.String("entity1"),
	}

	if err := ctx.DB.SaveEntity(e); err != nil {
		t.Fatal(err)
	}
}

func addGroup(t *testing.T, ctx tree.RefContext) {
	g := &pb.Group{
		Name:        proto.String("group1"),
		Number:      proto.Int32(1),
		DisplayName: proto.String("Group One"),
	}

	if err := ctx.DB.SaveGroup(g); err != nil {
		t.Fatal(err)
	}
}
