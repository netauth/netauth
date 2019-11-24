package interface_test

import (
	"testing"

	"github.com/golang/protobuf/proto"

	pb "github.com/netauth/protocol"
)

func TestBootstrap(t *testing.T) {
	em, ctx := newTreeManager(t)

	em.Bootstrap("entity1", "secret")

	e, err := ctx.DB.LoadEntity("entity1")
	if err != nil {
		t.Fatal(err)
	}

	if e.GetMeta().GetCapabilities()[0] != pb.Capability_GLOBAL_ROOT {
		t.Error("Entity wasn't bootstrapped")
	}

	// This only works because we're using nocrypto in the tests.
	if e.GetSecret() != "secret" {
		t.Error("Entity secret wasn't set")
	}
}

func TestBootstrapExistingEntity(t *testing.T) {
	em, ctx := newTreeManager(t)

	e := &pb.Entity{
		ID:     proto.String("_root"),
		Number: proto.Int32(1),
	}
	if err := ctx.DB.SaveEntity(e); err != nil {
		t.Fatal(err)
	}

	em.Bootstrap("_root", "secret")

	e, err := ctx.DB.LoadEntity("_root")
	if err != nil {
		t.Fatal(err)
	}

	if e.GetMeta().GetCapabilities()[0] != pb.Capability_GLOBAL_ROOT {
		t.Error("Entity wasn't bootstrapped")
	}

	// This only works because we're using nocrypto in the tests.
	if e.GetSecret() != "secret" {
		t.Error("Entity secret wasn't set")
	}
}
