package entity_manager

import "testing"

func TestNewEntityWithAuth(t *testing.T) {
	s := []struct {
		ID           string
		uidNumber    int32
		secret       string
		cap          string
		newID        string
		newUIDNumber int32
		newSecret    string
		wantErr      error
	}{
		{"a", -1, "a", "GLOBAL_ROOT", "aa", -1, "a", nil},
		{"b", -1, "a", "", "bb", -1, "a", E_ENTITY_UNQUALIFIED},
		{"c", -1, "a", "CREATE_ENTITY", "cc", -1, "a", nil},
	}

	resetMap()

	for _, c := range s {
		// Create entities with various capabilities
		if err := newEntity(c.ID, c.uidNumber, c.secret); err != nil {
			t.Error(err)
		}

		// Assign the test user some capabilities.
		if err := setEntityCapabilityByID(c.ID, c.cap); err != nil {
			t.Error(err)
		}

		// Test if those entities can perform various actions.
		if err := NewEntity(c.ID, c.secret, c.newID, c.newUIDNumber, c.newSecret); err != c.wantErr {
			t.Error(err)
		}
	}
}

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
		if err := newEntity(c.ID, c.uidNumber, c.secret); err != c.err {
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
		if err := newEntity(c.ID, c.uidNumber, c.secret); err != c.err {
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
		if err := newEntity(c.ID, c.uidNumber, c.secret); err != nil {
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
		if err := newEntity(c.ID, c.uidNumber, c.secret); err != nil {
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
		if err := newEntity(c.ID, c.uidNumber, c.secret); err != nil {
			t.Error(err)
		}

		if _, err := getEntityByID(c.ID); err != nil {
			t.Error("Added entity does not exist!")
		}
	}

	if _, err := getEntityByID("baz"); err == nil {
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
		if err := newEntity(c.ID, c.uidNumber, c.secret); err != nil {
			t.Error(err)
		}

		if _, err := getEntityByUIDNumber(c.uidNumber); err != nil {
			t.Error("Added entity does not exist!")
		}
	}

	if _, err := getEntityByUIDNumber(3); err == nil {
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
		newEntity(c.ID, c.uidNumber, c.secret)
		a, err := getEntityByID(c.ID)
		if err != nil {
			t.Error("Couldn't recall newly added entity!")
		}

		b, err := getEntityByUIDNumber(c.uidNumber)
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
		if err := newEntity(c.ID, c.uidNumber, c.secret); err != nil {
			t.Error(err)
		}
	}

	for _, c := range s {
		// Delete the entity
		if err := deleteEntityByID(c.ID); err != nil {
			t.Error(err)
		}

		// Make sure checking for that entity returns E_NO_ENTITY
		if _, err := getEntityByID(c.ID); err != E_NO_ENTITY {
			t.Error(err)
		}

		// Make sure that it is actually gone, and not just
		// gone from one index...
		if _, err := getEntityByUIDNumber(c.uidNumber); err != E_NO_ENTITY {
			t.Error(err)
		}
	}
}

func TestDeleteEntityWithAuth(t *testing.T) {
	x := []struct {
		ID        string
		uidNumber int32
		secret    string
	}{
		{"foo", 1, ""},
		{"bar", 2, ""},
		{"baz", 3, ""},
		{"quu", 4, ""},
	}

	s := []struct {
		ID        string
		uidNumber int32
		secret    string
		cap       string
		toDelete  string
		wantErr   error
	}{
		{"a", -1, "a", "GLOBAL_ROOT", "foo", nil},
		{"b", -1, "a", "", "bar", E_ENTITY_UNQUALIFIED},
		{"c", -1, "a", "CREATE_ENTITY", "baz", E_ENTITY_UNQUALIFIED},
		{"d", -1, "a", "DELETE_ENTITY", "quu", nil},
		{"e", -1, "a", "DELETE_ENTITY", "e", nil},
	}

	resetMap()

	for _, c := range x {
		// Create some entities to try deleting
		if err := newEntity(c.ID, c.uidNumber, c.secret); err != nil {
			t.Error(err)
		}
	}

	for _, c := range s {
		// Create entities with various capabilities
		if err := newEntity(c.ID, c.uidNumber, c.secret); err != nil {
			t.Error(err)
		}

		// Assign the test user some capabilities.
		if err := setEntityCapabilityByID(c.ID, c.cap); err != nil {
			t.Error(err)
		}

		// See if the entity can delete its target
		if err := DeleteEntity(c.ID, c.secret, c.toDelete); err != c.wantErr {
			t.Error(err)
		}
	}
}
