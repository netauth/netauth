package entity_manager

import "testing"

func TestAddDuplicateID(t *testing.T) {
	s := []struct {
		ID        string
		uidNumber int32
		secret    string
		err       error
	}{
		{"foo", 1, "", nil},
		{"foo", 2, "", E_DUPLICATE_ID},
	}

	// Force the map to be empty
	resetMap()

	for _, c := range s {
		if err := NewEntity(c.ID, c.uidNumber, c.secret); err != c.err {
			t.Errorf("Got %v; Want: %v", err, c.err)
		}
	}
}

func TestAddDuplicateUIDNumber(t *testing.T) {
	s := []struct {
		ID        string
		uidNumber int32
		secret    string
		err       error
	}{
		{"foo", 1, "", nil},
		{"bar", 1, "", E_DUPLICATE_UIDNUMBER},
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

func TestGetEntityByID(t *testing.T) {
	s := []struct {
		ID        string
		uidNumber int32
		secret    string
	}{
		{"foo", 1, ""},
		{"bar", 2, ""},
	}

	// Force the map to be empty
	resetMap()

	for _, c := range s {
		if err := NewEntity(c.ID, c.uidNumber, c.secret); err != nil {
			t.Error(err)
		}

		if _, err := GetEntityByID(c.ID); err != nil {
			t.Error("Added entity does not exist!")
		}
	}

	if _, err := GetEntityByID("baz"); err == nil {
		t.Error("Returned non-existant entity!")
	}
}

func TestGetEntityByUIDNumber(t *testing.T) {
	s := []struct {
		ID        string
		uidNumber int32
		secret    string
	}{
		{"foo", 1, ""},
		{"bar", 2, ""},
	}

	// Force the map to be empty
	resetMap()

	for _, c := range s {
		if err := NewEntity(c.ID, c.uidNumber, c.secret); err != nil {
			t.Error(err)
		}

		if _, err := GetEntityByUIDNumber(c.uidNumber); err != nil {
			t.Error("Added entity does not exist!")
		}
	}

	if _, err := GetEntityByUIDNumber(3); err == nil {
		t.Error("Returned non-existant entity!")
	}
}

func TestSameEntity(t *testing.T) {
	s := []struct {
		ID        string
		uidNumber int32
		secret    string
	}{
		{"foo", 1, ""},
		{"bar", 2, ""},
		{"baz", 3, ""},
	}

	resetMap()

	for _, c := range s {
		NewEntity(c.ID, c.uidNumber, c.secret)
		a, err := GetEntityByID(c.ID)
		if err != nil {
			t.Error("Couldn't recall newly added entity!")
		}

		b, err := GetEntityByUIDNumber(c.uidNumber)
		if err != nil {
			t.Error("Couldn't recall newly added entity!")
		}

		if a != b {
			t.Error("Different entities for same ID/Number!")
		}
	}
}

func TestDeleteEntityByID(t *testing.T) {
	s := []struct {
		ID        string
		uidNumber int32
		secret    string
	}{
		{"foo", 1, ""},
		{"bar", 2, ""},
		{"baz", 3, ""},
	}

	resetMap()

	// Populate some entities
	for _, c := range s {
		if err := NewEntity(c.ID, c.uidNumber, c.secret); err != nil {
			t.Error(err)
		}
	}

	for _, c := range s {
		// Delete the entity
		if err := DeleteEntityByID(c.ID); err != nil {
			t.Error(err)
		}

		// Make sure checking for that entity returns E_NO_ENTITY
		if _, err := GetEntityByID(c.ID); err != E_NO_ENTITY {
			t.Error(err)
		}

		// Make sure that it is actually gone, and not just
		// gone from one index...
		if _, err := GetEntityByUIDNumber(c.uidNumber); err != E_NO_ENTITY {
			t.Error(err)
		}
	}
}

func TestSetEntitySecretByID(t *testing.T) {
	s := []struct {
		ID        string
		uidNumber int32
		secret    string
	}{
		{"foo", 1, "a"},
		{"bar", 2, "a"},
		{"baz", 3, "a"},
	}

	resetMap()

	// Load in the entities
	for _, c := range s {
		if err := NewEntity(c.ID, c.uidNumber, c.secret); err != nil {
			t.Error(err)
		}
	}

	// Validate the secrets
	for _, c := range s {
		if err := ValidateEntitySecretByID(c.ID, c.secret); err != nil {
			t.Errorf("Failed: want 'nil', got %v", err)
		}
	}
}
