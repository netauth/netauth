package tree

import (
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/NetAuth/NetAuth/internal/db/impl/MemDB"
	"github.com/NetAuth/NetAuth/internal/crypto/impl/nocrypto"
	"github.com/NetAuth/NetAuth/pkg/errors"
	pb "github.com/NetAuth/Protocol"
)

func TestNewGroup(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

	s := []struct {
		name        string
		displayName string
		gidNumber   int32
		wantErr     error
	}{
		{"fooGroup", "", 1, nil},
		{"fooGroup", "", 1, errors.E_DUPLICATE_GROUP_ID},
		{"barGroup", "", 0, errors.E_DUPLICATE_GROUP_NUMBER},
		{"barGroup", "", 1, errors.E_DUPLICATE_GROUP_NUMBER},
		{"barGroup", "", -1, nil},
	}
	for _, c := range s {
		if err := em.NewGroup(c.name, c.displayName, "", c.gidNumber); err != c.wantErr {
			t.Errorf("Wrong Error: want '%v' got '%v'", c.wantErr, err)
		}
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

func TestDeleteGroup(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

	if err := em.NewGroup("foo", "", "", -1); err != nil {
		t.Error(err)
	}

	if _, err := em.GetGroupByName("foo"); err != nil {
		t.Error(err)
	}

	if err := em.DeleteGroup("foo"); err != nil {
		t.Error(err)
	}

	if _, err := em.GetGroupByName("foo"); err != errors.E_NO_GROUP {
		t.Error(err)
	}
}

func TestUpdateGroupMetaInternal(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

	if err := em.NewGroup("foo", "foo", "", -1); err != nil {
		t.Error(err)
	}

	update := &pb.Group{DisplayName: proto.String("Foo Group")}

	if err := em.UpdateGroupMeta("foo", update); err != nil {
		t.Error(err)
	}

	g, err := em.GetGroupByName("foo")
	if err != nil {
		t.Error(err)
	}

	if g.GetDisplayName() != "Foo Group" {
		t.Error("Meta update failed!")
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
