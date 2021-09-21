package interface_test

import (
	"context"
	"testing"

	"google.golang.org/protobuf/proto"

	pb "github.com/netauth/protocol"
)

func TestRemoveEntityFromGroup(t *testing.T) {
	ctxt := context.Background()
	m, ctx := newTreeManager(t)

	e := &pb.Entity{
		ID: proto.String("entity1"),
		Meta: &pb.EntityMeta{
			Groups: []string{
				"group1",
				"group2",
			},
		},
	}

	if err := ctx.DB.SaveEntity(ctxt, e); err != nil {
		t.Fatal(err)
	}

	if err := m.RemoveEntityFromGroup(ctxt, "entity1", "group1"); err != nil {
		t.Fatal(err)
	}

	e, err := ctx.DB.LoadEntity(ctxt, "entity1")
	if err != nil {
		t.Fatal(err)
	}

	groups := e.GetMeta().GetGroups()
	if len(groups) != 1 || groups[0] != "group2" {
		t.Error("Entity modification error")
	}
}
