package interface_test

import (
	"sort"
	"testing"

	"github.com/golang/protobuf/proto"

	pb "github.com/netauth/protocol"
)

func TestGetMemberships(t *testing.T) {
	m, ctx := newTreeManager(t)

	buildSampleTree(t, ctx)

	// entity1, indirects off
	e, err := ctx.DB.LoadEntity("entity1")
	if err != nil {
		t.Fatal(err)
	}
	result := m.GetMemberships(e, false)
	sort.Strings(result)
	if len(result) != 2 || result[0] != "group1" || result[1] != "group2" {
		t.Log(result)
		t.Error("entity1 has the wrong group membership")
	}

	// entity1, indirects on
	result = m.GetMemberships(e, true)
	sort.Strings(result)
	if len(result) != 3 || result[0] != "group1" || result[1] != "group2" || result[2] != "group4" {
		t.Log(result)
		t.Error("entity1 has the wrong group membership")
	}
}

func TestGetMembershipsBadEntity(t *testing.T) {
	m, _ := newTreeManager(t)

	e := &pb.Entity{
		ID: proto.String("unknown"),
	}

	if result := m.GetMemberships(e, true); len(result) != 0 {
		t.Fatal("Memberships found for an unknown group")
	}
}
