package entity_manager

import (
	"testing"

	"github.com/NetAuth/NetAuth/internal/server/crypto/impl/nocrypto"
	"github.com/NetAuth/NetAuth/internal/server/db/impl/MemDB"

	pb "github.com/NetAuth/NetAuth/pkg/proto"
)

func TestInternalDirectMembershipEdit(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

	e := &pb.Entity{}

	if err := em.addEntityToDirectGroup(e, "fooGroup"); err != nil {
		t.Error(err)
	}

	groups := em.getDirectGroups(e)
	if len(groups) != 1 || groups[0] != "fooGroup" {
		t.Error("Wrong group number/membership")
	}

	em.removeEntityFromDirectGroup(e, "fooGroup")
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
	em.removeEntityFromDirectGroup(e, "fooGroup")
}

func TestGetDirectGroupsNoMeta(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

	e := &pb.Entity{}

	if groups := em.getDirectGroups(e); len(groups) !=0 {
		t.Error("getDirectGroups fabricated a group!")
	}
}
