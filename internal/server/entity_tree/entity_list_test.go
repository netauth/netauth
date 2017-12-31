package entity_tree

import (
	"testing"

	pb "github.com/NetAuth/NetAuth/proto"
)

func resetMap() {
	e = make(map[string]*pb.Entity)
}

func TestAddFirstEntity(t *testing.T) {
	ID := "foo"
	uidNumber := int32(1)
	secret := "secret"

	// Force the map to be empty
	resetMap()

	// Create an entity
	if err := NewEntity(ID, uidNumber, secret); err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	// Verify that the map now contains exactly the entity specified
	if e[ID].GetID() != ID {
		t.Errorf("Expected uid to be %v got %v", e[ID].GetID(), ID)
	}
	if e[ID].GetUidNumber() != uidNumber {
		t.Errorf("Expected uidNumber to be %v got %v", e[ID].GetUidNumber(), uidNumber)
	}
	if e[ID].GetSecret() != secret {
		t.Errorf("Expected secret to be %v got %v", e[ID].GetSecret(), secret)
	}
}

func TestAddDuplicateEntity(t *testing.T) {
	s := []struct {
		ID        string
		uidNumber int32
		secret    string
		err       error
	}{
		{"foo", 1, "bar", nil},
		{"foo", 1, "bar", E_DUPLICATE_ID},
	}

	// Force the map to be empty
	resetMap()

	for _, c := range s {
		if err := NewEntity(c.ID, c.uidNumber, c.secret); err != c.err {
			t.Errorf("Got %v; Want: %v", err, c.err)
		}
	}
}

func TestNextUIDNumber(t *testing.T) {
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

	resetMap()

	for _, c := range s {
		//  Make sure the entity actually gets added
		if err := NewEntity(c.ID, c.uidNumber, c.secret); err != nil {
			t.Error(err)
		}

		// Validate that after a given mutation the number is
		// still what we expect it to be.
		if next := nextUIDNumber(); next != c.nextUIDNumber {
			t.Errorf("Wrong next number; got: %v want %v", next, c.nextUIDNumber)
		}
	}

}

func TestNewEntityAutoNumber(t *testing.T) {
	s := []struct {
		ID        string
		uidNumber int32
		secret    string
	}{
		{"foo", 1, ""},
		{"bar", -1, ""},
		{"baz", 3, ""},
	}

	// Force the map to be empty
	resetMap()

	for _, c := range s {
		if err := NewEntity(c.ID, c.uidNumber, c.secret); err != nil {
			t.Error(err)
		}
	}
}
