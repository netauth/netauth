package interface_test

import (
	"context"
	"testing"

	"google.golang.org/protobuf/proto"

	pb "github.com/netauth/protocol"
)

func TestModifyGroupExpansions(t *testing.T) {
	ctxt := context.Background()
	m, mdb := newTreeManager(t)

	g1 := &pb.Group{
		Name: proto.String("group1"),
	}
	g2 := &pb.Group{
		Name: proto.String("group2"),
	}

	if err := mdb.SaveGroup(ctxt, g1); err != nil {
		t.Fatal(err)
	}
	if err := mdb.SaveGroup(ctxt, g2); err != nil {
		t.Fatal(err)
	}

	if err := m.ModifyGroupExpansions(ctxt, "group1", "group2", pb.ExpansionMode_INCLUDE); err != nil {
		t.Fatal(err)
	}

	g, err := mdb.LoadGroup(ctxt, "group1")
	if err != nil {
		t.Fatal(err)
	}

	if len(g.GetExpansions()) != 1 || g.GetExpansions()[0] != "INCLUDE:group2" {
		t.Error("Expansions are not correctly set on group1")
	}
}
