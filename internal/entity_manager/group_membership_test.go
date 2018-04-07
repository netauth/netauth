package entity_manager

import (
	"testing"

	"github.com/NetAuth/NetAuth/internal/crypto/impl/nocrypto"
	"github.com/NetAuth/NetAuth/internal/db/impl/MemDB"

	pb "github.com/NetAuth/NetAuth/pkg/proto"
)

func TestMembershipEdit(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

	e := &pb.Entity{}

	if err := em.NewGroup("fooGroup", "", 1000); err != nil {
		t.Error(err)
	}

	if err := em.AddEntityToGroup(e, "fooGroup"); err != nil {
		t.Error(err)
	}

	groups := em.GetDirectGroups(e)
	if len(groups) != 1 || groups[0] != "fooGroup" {
		t.Error("Wrong group number/membership")
	}

	em.RemoveEntityFromGroup(e, "fooGroup")
	groups = em.GetDirectGroups(e)
	if len(groups) != 0 {
		t.Error("Wrong group number/membership")
	}
}

func TestRemoveEntityFromGroupNilMeta(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

	e := &pb.Entity{}

	// This is just to make sure that this function doesn't
	// explode.
	em.RemoveEntityFromGroup(e, "fooGroup")
}

func TestGetGroupsNoMeta(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

	e := &pb.Entity{}

	if groups := em.GetDirectGroups(e); len(groups) != 0 {
		t.Error("getDirectGroups fabricated a group!")
	}
}
