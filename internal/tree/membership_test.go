package tree

import (
	"testing"

	"github.com/NetAuth/NetAuth/internal/crypto/impl/nocrypto"
	"github.com/NetAuth/NetAuth/internal/db/impl/MemDB"
	"github.com/NetAuth/NetAuth/pkg/errors"

	pb "github.com/NetAuth/Protocol"
)

func TestMembershipEdit(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

	e := &pb.Entity{}

	if err := em.NewGroup("fooGroup", "", "", 1000); err != nil {
		t.Error(err)
	}

	if err := em.addEntityToGroup(e, "fooGroup"); err != nil {
		t.Error(err)
	}

	groups := em.GetDirectGroups(e)
	if len(groups) != 1 || groups[0] != "fooGroup" {
		t.Error("Wrong group number/membership")
	}

	em.removeEntityFromGroup(e, "fooGroup")
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
	em.removeEntityFromGroup(e, "fooGroup")
}

func TestGetGroupsNoMeta(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

	e := &pb.Entity{}

	if groups := em.GetDirectGroups(e); len(groups) != 0 {
		t.Error("getDirectGroups fabricated a group!")
	}
}

func TestListMembersALLInternal(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

	s := []string{
		"foo",
		"bar",
		"baz",
	}

	for _, id := range s {
		if err := em.NewEntity(id, -1, ""); err != nil {
			t.Error(err)
		}

		listALL, err := em.listMembers("ALL")
		if err != nil {
			t.Error(err)
		}

		dbAll, err := em.db.DiscoverEntityIDs()
		if err != nil {
			t.Error(err)
		}
		if len(dbAll) != len(listALL) {
			t.Error("Different number of entities returned!")
		}
	}
}

func TestListMembersNoMatchInternal(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())
	list, err := em.listMembers("")
	if list != nil && err != errors.E_NO_GROUP {
		t.Error(err)
	}
}

func TestListMembersExternal(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

	s := []string{
		"foo",
		"bar",
		"baz",
	}

	for _, id := range s {
		if err := em.NewEntity(id, -1, ""); err != nil {
			t.Error(err)
		}

		listALLExt, err := em.ListMembers("ALL")
		if err != nil {
			t.Error(err)
		}

		listALLInt, err := em.listMembers("ALL")
		if err != nil {
			t.Error(err)
		}

		// At first this doesn't look like it tests anything,
		// but actually what's being tested here is that the
		// same number of elements comes back from both the
		// internal version and the copy version.  The copy
		// itself is tested in another case.
		if len(listALLExt) != len(listALLInt) {
			t.Error("Different sizes for same group!?")
		}
	}
}

func TestListMembersNoMatchExternal(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())
	list, err := em.ListMembers("")
	if list != nil && err != errors.E_NO_GROUP {
		t.Error(err)
	}
}
