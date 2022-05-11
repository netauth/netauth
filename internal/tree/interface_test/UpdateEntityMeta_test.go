package interface_test

import (
	"context"
	"testing"

	"google.golang.org/protobuf/proto"

	pb "github.com/netauth/protocol"
)

func TestUpdateEntityMeta(t *testing.T) {
	ctxt := context.Background()
	m, mdb := newTreeManager(t)

	addEntity(t, mdb)

	meta := &pb.EntityMeta{
		GECOS: proto.String("A Test Entity"),
	}

	if err := m.UpdateEntityMeta(ctxt, "entity1", meta); err != nil {
		t.Fatal(err)
	}

	e, err := mdb.LoadEntity(ctxt, "entity1")
	if err != nil {
		t.Fatal(err)
	}

	if e.GetMeta().GetGECOS() != "A Test Entity" {
		t.Error("Metadata not set")
	}
}
