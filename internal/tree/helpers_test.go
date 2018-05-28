package tree

import (
	"testing"

	"github.com/NetAuth/NetAuth/internal/db/impl/MemDB"
	"github.com/NetAuth/NetAuth/internal/crypto/impl/nocrypto"
	"github.com/golang/protobuf/proto"
)

func TestNextUIDNumber(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

	s := []struct {
		ID            string
		number     int32
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
		if err := em.NewEntity(c.ID, c.number, c.secret); err != nil {
			t.Error(err)
		}

		// Validate that after a given mutation the number is
		// still what we expect it to be.
		next, err := em.nextUIDNumber()
		if err != nil {
			t.Error(err)
		}
		if next != c.nextUIDNumber {
			t.Errorf("Wrong next number; got: %v want %v", next, c.nextUIDNumber)
		}
	}

}

func TestGetEntityByID(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

	s := []struct {
		ID        string
		number int32
		secret    string
	}{
		{"foo", 1, ""},
		{"bar", 2, ""},
	}

	for _, c := range s {
		if err := em.NewEntity(c.ID, c.number, c.secret); err != nil {
			t.Error(err)
		}

		if _, err := em.db.LoadEntity(c.ID); err != nil {
			t.Error("Added entity does not exist!")
		}
	}

	if _, err := em.db.LoadEntity("baz"); err == nil {
		t.Error("Returned non-existant entity!")
	}
}

func TestSafeCopyEntity(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

	if err := em.NewEntity("foo", -1, "bar"); err != nil {
		t.Error(err)
	}

	e, err := em.db.LoadEntity("foo")
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
