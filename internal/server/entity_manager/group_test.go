package entity_manager

import (
	"testing"
)

func TestListMembersALLInternal(t *testing.T) {
	em := New()

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

		if len(em.eByID) != len(listALL) {
			t.Error("Different number of entities returned!")
		}
	}
}

func TestListMembersNoMatchInternal(t *testing.T) {
	em := New()
	list, err := em.listMembers("")
	if list != nil && err != E_NO_GROUP {
		t.Error(err)
	}
}

func TestListMembersExternal(t *testing.T) {
	em := New()

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
	em := New()
	list, err := em.ListMembers("")
	if list != nil && err != E_NO_GROUP {
		t.Error(err)
	}
}
