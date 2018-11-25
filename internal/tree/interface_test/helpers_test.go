package interface_test

import (
	"testing"

	"github.com/golang/protobuf/proto"
	
	"github.com/NetAuth/NetAuth/internal/crypto/nocrypto"
	"github.com/NetAuth/NetAuth/internal/db/memdb"
	"github.com/NetAuth/NetAuth/internal/tree"
	_ "github.com/NetAuth/NetAuth/internal/tree/hooks"

	pb "github.com/NetAuth/Protocol"
)

func newTreeManager(t *testing.T) (*tree.Manager, tree.RefContext) {
	mdb, err := memdb.New()
	if err != nil {
		t.Fatal(err)
	}

	crypto, err := nocrypto.New()
	if err != nil {
		t.Fatal(err)
	}

	ctx := tree.RefContext{
		DB: mdb,
		Crypto: crypto,
	}
	
	em := tree.New(ctx.DB, ctx.Crypto)

	return em, ctx
}

func addEntity(t *testing.T, ctx tree.RefContext) {
	e := &pb.Entity{
		ID: proto.String("entity1"),
		Number: proto.Int32(1),
		Secret: proto.String("entity1"),
	}

	if err := ctx.DB.SaveEntity(e); err != nil {
		t.Fatal(err)
	}
}
