package entity_manager

import (
	"testing"

	"github.com/NetAuth/NetAuth/internal/server/db/impl/MemDB"
	"github.com/NetAuth/NetAuth/pkg/errors"
)

func TestNewGroupInternal(t *testing.T) {
	em := New(MemDB.New())

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
		if err := em.newGroup(c.name, c.displayName, c.gidNumber); err != c.wantErr {
			t.Errorf("Wrong Error: want '%v' got '%v'", c.wantErr, err)
		}
	}
}

func TestGroupExternal(t *testing.T) {
	em := New(MemDB.New())

	// Add some users to the system that can and cannot create
	// groups.
	e := []struct {
		ID     string
		secret string
		cap    string
	}{
		{"foo", "foo", ""},
		{"bar", "bar", "CREATE_GROUP"},
	}
	for _, ne := range e {
		if err := em.newEntity(ne.ID, -1, ne.secret); err != nil {
			t.Error(err)
		}

		if err := em.setEntityCapabilityByID(ne.ID, ne.cap); err != nil {
			t.Error(err)
		}
	}

	s := []struct {
		ID        string
		secret    string
		name      string
		gidNumber int32
		wantErr   error
	}{
		{"foo", "foo", "newGroup", -1, errors.E_ENTITY_UNQUALIFIED},
		{"bar", "bar", "newGroup", -1, nil},
		{"bar", "bar", "newGroup", -1, errors.E_DUPLICATE_GROUP_ID},
		{"bar", "bar", "fooGroup", 0, errors.E_DUPLICATE_GROUP_NUMBER},
	}

	for _, c := range s {
		if err := em.NewGroup(c.ID, c.secret, c.name, "", c.gidNumber); err != c.wantErr {
			t.Error(err)
		}
	}
}

func TestListMembersALLInternal(t *testing.T) {
	em := New(MemDB.New())

	s := []string{
		"foo",
		"bar",
		"baz",
	}

	for _, id := range s {
		if err := em.newEntity(id, -1, ""); err != nil {
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
	em := New(MemDB.New())
	list, err := em.listMembers("")
	if list != nil && err != errors.E_NO_GROUP {
		t.Error(err)
	}
}

func TestListMembersExternal(t *testing.T) {
	em := New(MemDB.New())

	s := []string{
		"foo",
		"bar",
		"baz",
	}

	for _, id := range s {
		if err := em.newEntity(id, -1, ""); err != nil {
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
	em := New(MemDB.New())
	list, err := em.ListMembers("")
	if list != nil && err != errors.E_NO_GROUP {
		t.Error(err)
	}
}
