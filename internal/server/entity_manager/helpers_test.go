package entity_manager

import (
	"testing"

	"github.com/golang/protobuf/proto"
)

func TestNextUIDNumber(t *testing.T) {
	em := New(nil)
	em.initMem()

	s := []struct {
		ID            string
		uidNumber     int32
		secret        string
		nextUIDNumber int32
	}{
		{"foo", 1, "", 2},
		{"bar", 2, "", 3},
		{"baz", 65, "", 66}, // Numbers may be missing in the middle
		{"fuu", 23, "", 66}, // Later additions shouldn't alter max
	}

	for _, c := range s {
		//  Make sure the entity actually gets added
		if err := em.newEntity(c.ID, c.uidNumber, c.secret); err != nil {
			t.Error(err)
		}

		// Validate that after a given mutation the number is
		// still what we expect it to be.
		if next := em.nextUIDNumber(); next != c.nextUIDNumber {
			t.Errorf("Wrong next number; got: %v want %v", next, c.nextUIDNumber)
		}
	}

}

func TestGetEntityByID(t *testing.T) {
	em := New(nil)
	em.initMem()

	s := []struct {
		ID        string
		uidNumber int32
		secret    string
	}{
		{"foo", 1, ""},
		{"bar", 2, ""},
	}

	for _, c := range s {
		if err := em.newEntity(c.ID, c.uidNumber, c.secret); err != nil {
			t.Error(err)
		}

		if _, err := em.getEntityByID(c.ID); err != nil {
			t.Error("Added entity does not exist!")
		}
	}

	if _, err := em.getEntityByID("baz"); err == nil {
		t.Error("Returned non-existant entity!")
	}
}

func TestSafeCopyEntity(t *testing.T) {
	em := New(nil)
	em.initMem()

	if err := em.newEntity("foo", -1, "bar"); err != nil {
		t.Error(err)
	}

	e, err := em.getEntityByID("foo")
	if err != nil {
		t.Error(err)
	}

	ne, err := safeCopyEntity(e)
	if err != nil {
		t.Error(err)
	}

	// The normal way to do this would be to check if the proto is
	// the same, but here we need to check if two fields are
	// different, then make sure that everything else is the same.
	if e.GetSecret() == ne.GetSecret() {
		t.Error("Secret field not obscured!")
	}

	e.Secret = proto.String("")
	ne.Secret = proto.String("")

	if !proto.Equal(e, ne) {
		t.Error("Entity values not otherwise equal!")
	}
}
