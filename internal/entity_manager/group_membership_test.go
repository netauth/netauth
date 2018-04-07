package entity_manager

import (
	"testing"

	"github.com/NetAuth/NetAuth/internal/crypto/impl/nocrypto"
	"github.com/NetAuth/NetAuth/internal/db/impl/MemDB"

	pb "github.com/NetAuth/NetAuth/pkg/proto"
)

func TestInternalMembershipEdit(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

	e := &pb.Entity{}

	if err := em.newGroup("fooGroup", "", 1000); err != nil {
		t.Error(err)
	}

	if err := em.addEntityToGroup(e, "fooGroup"); err != nil {
		t.Error(err)
	}

	groups := em.getDirectGroups(e)
	if len(groups) != 1 || groups[0] != "fooGroup" {
		t.Error("Wrong group number/membership")
	}

	em.removeEntityFromGroup(e, "fooGroup")
	groups = em.getDirectGroups(e)
	if len(groups) != 0 {
		t.Error("Wrong group number/membership")
	}
}

func TestRemoveEntityFromGroupInternalNilMeta(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

	e := &pb.Entity{}

	// This is just to make sure that this function doesn't
	// explode.
	em.removeEntityFromGroup(e, "fooGroup")
}

func TestGetGroupsNoMeta(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

	e := &pb.Entity{}

	if groups := em.getDirectGroups(e); len(groups) != 0 {
		t.Error("getDirectGroups fabricated a group!")
	}
}
