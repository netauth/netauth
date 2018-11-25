package interface_test

import (
	"testing"

	"github.com/golang/protobuf/proto"

	pb "github.com/NetAuth/Protocol"
)

func TestUpdateEntityMeta(t *testing.T) {
	m, ctx := newTreeManager(t)

	addEntity(t, ctx)

	meta := &pb.EntityMeta{
		GECOS: proto.String("A Test Entity"),
	}

	if err := m.UpdateEntityMeta("entity1", meta); err != nil {
		t.Fatal(err)
	}

	e, err := ctx.DB.LoadEntity("entity1")
	if err != nil {
		t.Fatal(err)
	}

	if e.GetMeta().GetGECOS() != "A Test Entity" {
		t.Error("Metadata not set")
	}
}
