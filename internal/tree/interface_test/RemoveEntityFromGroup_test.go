package interface_test

import (
	"testing"

	"google.golang.org/protobuf/proto"

	pb "github.com/netauth/protocol"
)

func TestRemoveEntityFromGroup(t *testing.T) {
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

	if err := ctx.DB.SaveEntity(e); err != nil {
		t.Fatal(err)
	}

	if err := m.RemoveEntityFromGroup("entity1", "group1"); err != nil {
		t.Fatal(err)
	}

	e, err := ctx.DB.LoadEntity("entity1")
	if err != nil {
		t.Fatal(err)
	}

	groups := e.GetMeta().GetGroups()
	if len(groups) != 1 || groups[0] != "group2" {
		t.Error("Entity modification error")
	}
}
