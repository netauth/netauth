package interface_test

import (
	"context"
	"sort"
	"testing"

	"google.golang.org/protobuf/proto"

	"github.com/netauth/netauth/internal/db"

	pb "github.com/netauth/protocol"
)

func TestGetMemberships(t *testing.T) {
	ctxt := context.Background()
	m, ctx := newTreeManager(t)

	buildSampleTree(t, ctx)
	ctx.DB.(*db.DB).EventUpdateAll()

	// entity1, indirects off
	e, err := ctx.DB.LoadEntity(ctxt, "entity1")
	if err != nil {
		t.Fatal(err)
	}
	result := m.GetMemberships(ctxt, e)
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

	if result := m.GetMemberships(context.Background(), e); len(result) != 0 {
		t.Fatal("Memberships found for an unknown group")
	}
}
