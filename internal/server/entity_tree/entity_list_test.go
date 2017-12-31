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
	s := []struct{
		ID string
		uidNumber int32
		secret string
		err error
	}{
		{"foo", 1, "bar", nil},
		{"foo", 1, "bar", E_DUPLICATE_ID},
	}
	
	// Force the map to be empty
	resetMap()

	for _, c := range(s) {
		err := NewEntity(c.ID, c.uidNumber, c.secret)
		if err != c.err {
			t.Errorf("Got %v; Want: %v", err, c.err)
		}
	}
}
