package interface_test

import (
	"testing"

	"github.com/golang/protobuf/proto"

	pb "github.com/netauth/protocol"
	rpc "github.com/netauth/protocol/v2"
)

func TestModifyGroupRule(t *testing.T) {
	m, ctx := newTreeManager(t)

	g1 := &pb.Group{
		Name: proto.String("group1"),
	}
	g2 := &pb.Group{
		Name: proto.String("group2"),
	}

	if err := ctx.DB.SaveGroup(g1); err != nil {
		t.Fatal(err)
	}
	if err := ctx.DB.SaveGroup(g2); err != nil {
		t.Fatal(err)
	}

	if err := m.ModifyGroupRule("group1", "group2", rpc.RuleAction_INCLUDE); err != nil {
		t.Fatal(err)
	}

	g, err := ctx.DB.LoadGroup("group1")
	if err != nil {
		t.Fatal(err)
	}

	if len(g.GetExpansions()) != 1 || g.GetExpansions()[0] != "INCLUDE:group2" {
		t.Error(g.GetExpansions())
		t.Error("Expansions are not correctly set on group1")
	}

	if err := m.ModifyGroupRule("group1", "group2", rpc.RuleAction_REMOVE_RULE); err != nil {
		t.Fatal(err)
	}

	g, err = ctx.DB.LoadGroup("group1")
	if err != nil {
		t.Fatal(err)
	}

	if len(g.GetExpansions()) != 0 {
		t.Error(g.GetExpansions())
		t.Error("Expansions are not correctly set on group1")
	}

	if err := m.ModifyGroupRule("group1", "group2", rpc.RuleAction_EXCLUDE); err != nil {
		t.Fatal(err)
	}

	g, err = ctx.DB.LoadGroup("group1")
	if err != nil {
		t.Fatal(err)
	}

	if len(g.GetExpansions()) != 1 || g.GetExpansions()[0] != "EXCLUDE:group2" {
		t.Error(g.GetExpansions())
		t.Error("Expansions are not correctly set on group1")
	}
}
