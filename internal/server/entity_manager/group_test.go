package entity_manager

import (
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/NetAuth/NetAuth/internal/server/db/impl/MemDB"
	"github.com/NetAuth/NetAuth/internal/server/crypto/impl/nocrypto"
	"github.com/NetAuth/NetAuth/pkg/errors"
	pb "github.com/NetAuth/NetAuth/pkg/proto"
)

func TestNewGroupInternal(t *testing.T) {
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
		if err := em.newGroup(c.name, c.displayName, c.gidNumber); err != c.wantErr {
			t.Errorf("Wrong Error: want '%v' got '%v'", c.wantErr, err)
		}
	}
}

func TestNewGroupExternal(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

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
	em := New(MemDB.New(), nocrypto.New())

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

func TestDeleteGroupInternal(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

	if err := em.newGroup("foo", "", -1); err != nil {
		t.Error(err)
	}

	if _, err := em.getGroupByName("foo"); err != nil {
		t.Error(err)
	}

	if err := em.deleteGroup("foo"); err != nil {
		t.Error(err)
	}

	if _, err := em.getGroupByName("foo"); err != errors.E_NO_GROUP {
		t.Error(err)
	}
}

func TestDeleteGroupExternal(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

	// Add some users to the system that can and cannot delete
	// groups.
	e := []struct {
		ID     string
		secret string
		cap    string
	}{
		{"foo", "foo", ""},
		{"bar", "bar", "DESTROY_GROUP"},
	}
	for _, ne := range e {
		if err := em.newEntity(ne.ID, -1, ne.secret); err != nil {
			t.Error(err)
		}

		if err := em.setEntityCapabilityByID(ne.ID, ne.cap); err != nil {
			t.Error(err)
		}
	}

	if err := em.newGroup("foo", "", -1); err != nil {
		t.Error(err)
	}

	if err := em.DeleteGroup("foo", "foo", "foo"); err != errors.E_ENTITY_UNQUALIFIED {
		t.Error(err)
	}

	if err := em.DeleteGroup("bar", "bar", "foo"); err != nil {
		t.Error(err)
	}

	if err := em.DeleteGroup("bar", "bar", "foo"); err != errors.E_NO_GROUP {
		t.Error(err)
	}
}

func TestUpdateGroupMetaInternal(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

	if err := em.newGroup("foo", "foo", -1); err != nil {
		t.Error(err)
	}

	update := &pb.Group{DisplayName: proto.String("Foo Group")}

	if err := em.updateGroupMeta("foo", update); err != nil {
		t.Error(err)
	}

	g, err := em.getGroupByName("foo")
	if err != nil {
		t.Error(err)
	}

	if g.GetDisplayName() != "Foo Group" {
		t.Error("Meta update failed!")
	}
}

func TestUpdateGroupMetaExternal(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

	// Add some users to the system that can and cannot modify
	// groups.
	e := []struct {
		ID     string
		secret string
		cap    string
	}{
		{"foo", "foo", ""},
		{"bar", "bar", "MODIFY_GROUP_META"},
	}
	for _, ne := range e {
		if err := em.newEntity(ne.ID, -1, ne.secret); err != nil {
			t.Error(err)
		}

		if err := em.setEntityCapabilityByID(ne.ID, ne.cap); err != nil {
			t.Error(err)
		}
	}

	if err := em.newGroup("foo", "", -1); err != nil {
		t.Error(err)
	}

	update := &pb.Group{DisplayName: proto.String("foo display name")}

	if err := em.UpdateGroupMeta("foo", "foo", "foo", update); err != errors.E_ENTITY_UNQUALIFIED {
		t.Error(err)
	}

	if err := em.UpdateGroupMeta("bar", "bar", "foo", update); err != nil {
		t.Error(err)
	}

	g, err := em.getGroupByName("foo")
	if err != nil {
		t.Error(err)
	}
	if g.GetDisplayName() != "foo display name" {
		t.Error("Group Update failed!")
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
	em := New(MemDB.New(), nocrypto.New())
	list, err := em.ListMembers("")
	if list != nil && err != errors.E_NO_GROUP {
		t.Error(err)
	}
}
