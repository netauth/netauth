package entity_manager

import "testing"

func TestNewEntityWithAuth(t *testing.T) {
	em := New()

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

	for _, c := range s {
		// Create entities with various capabilities
		if err := em.newEntity(c.ID, c.uidNumber, c.secret); err != nil {
			t.Error(err)
		}

		// Assign the test user some capabilities.
		if err := em.setEntityCapabilityByID(c.ID, c.cap); err != nil {
			t.Error(err)
		}

		// Test if those entities can perform various actions.
		if err := em.NewEntity(c.ID, c.secret, c.newID, c.newUIDNumber, c.newSecret); err != c.wantErr {
			t.Error(err)
		}
	}
}

func TestAddDuplicateID(t *testing.T) {
	em := New()

	s := []struct {
		ID        string
		uidNumber int32
		secret    string
		err       error
	}{
		{"foo", 1, "", nil},
		{"foo", 2, "", E_DUPLICATE_ID},
	}

	for _, c := range s {
		if err := em.newEntity(c.ID, c.uidNumber, c.secret); err != c.err {
			t.Errorf("Got %v; Want: %v", err, c.err)
		}
	}
}

func TestAddDuplicateUIDNumber(t *testing.T) {
	em := New()

	s := []struct {
		ID        string
		uidNumber int32
		secret    string
		err       error
	}{
		{"foo", 1, "", nil},
		{"bar", 1, "", E_DUPLICATE_UIDNUMBER},
	}

	for _, c := range s {
		if err := em.newEntity(c.ID, c.uidNumber, c.secret); err != c.err {
			t.Errorf("Got %v; Want: %v", err, c.err)
		}
	}
}

func TestNextUIDNumber(t *testing.T) {
	em := New()

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

func TestNewEntityAutoNumber(t *testing.T) {
	em := New()

	s := []struct {
		ID        string
		uidNumber int32
		secret    string
	}{
		{"foo", 1, ""},
		{"bar", -1, ""},
		{"baz", 3, ""},
	}

	for _, c := range s {
		if err := em.newEntity(c.ID, c.uidNumber, c.secret); err != nil {
			t.Error(err)
		}
	}
}

func TestMakeBootstrap(t *testing.T) {
	em := New()

	s := []struct {
		ID            string
		secret        string
		pre_disable   bool
		bootstrap_val bool
		wantErr       error
	}{
		{"bar", "", false, false, nil},
		{"baz", "", false, false, nil},
		{"quu", "", true, true, E_NO_ENTITY},
		{"qux", "", true, false, E_NO_ENTITY},
	}

	// Precreate bar to make sure they can get the
	// global_superuser power in the existing path
	if err := em.newEntity("bar", -1, ""); err != nil {
		t.Error(err)
	}

	for _, c := range s {
		em.bootstrap_done = c.bootstrap_val
		if c.pre_disable {
			em.DisableBootstrap()
		}
		em.MakeBootstrap(c.ID, c.secret)

		if err := em.checkEntityCapabilityByID(c.ID, "GLOBAL_ROOT"); err != c.wantErr {
			t.Error(err)
		}
	}
}

func TestGetEntityByID(t *testing.T) {
	em := New()

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

func TestGetEntityByUIDNumber(t *testing.T) {
	em := New()

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

		if _, err := em.getEntityByUIDNumber(c.uidNumber); err != nil {
			t.Error("Added entity does not exist!")
		}
	}

	if _, err := em.getEntityByUIDNumber(3); err == nil {
		t.Error("Returned non-existant entity!")
	}
}

func TestSameEntity(t *testing.T) {
	em := New()

	s := []struct {
		ID        string
		uidNumber int32
		secret    string
	}{
		{"foo", 1, ""},
		{"bar", 2, ""},
		{"baz", 3, ""},
	}

	for _, c := range s {
		if err := em.newEntity(c.ID, c.uidNumber, c.secret); err != nil {
			t.Error(err)
		}

		a, err := em.getEntityByID(c.ID)
		if err != nil {
			t.Error("Couldn't recall newly added entity!")
		}

		b, err := em.getEntityByUIDNumber(c.uidNumber)
		if err != nil {
			t.Error("Couldn't recall newly added entity!")
		}

		if a != b {
			t.Error("Different entities for same ID/Number!")
		}
	}
}

func TestDeleteEntityByID(t *testing.T) {
	em := New()

	s := []struct {
		ID        string
		uidNumber int32
		secret    string
	}{
		{"foo", 1, ""},
		{"bar", 2, ""},
		{"baz", 3, ""},
	}

	// Populate some entities
	for _, c := range s {
		if err := em.newEntity(c.ID, c.uidNumber, c.secret); err != nil {
			t.Error(err)
		}
	}

	for _, c := range s {
		// Delete the entity
		if err := em.deleteEntityByID(c.ID); err != nil {
			t.Error(err)
		}

		// Make sure checking for that entity returns E_NO_ENTITY
		if _, err := em.getEntityByID(c.ID); err != E_NO_ENTITY {
			t.Error(err)
		}

		// Make sure that it is actually gone, and not just
		// gone from one index...
		if _, err := em.getEntityByUIDNumber(c.uidNumber); err != E_NO_ENTITY {
			t.Error(err)
		}
	}
}

func TestDeleteEntityWithAuth(t *testing.T) {
	em := New()

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

	for _, c := range x {
		// Create some entities to try deleting
		if err := em.newEntity(c.ID, c.uidNumber, c.secret); err != nil {
			t.Error(err)
		}
	}

	for _, c := range s {
		// Create entities with various capabilities
		if err := em.newEntity(c.ID, c.uidNumber, c.secret); err != nil {
			t.Error(err)
		}

		// Assign the test user some capabilities.
		if err := em.setEntityCapabilityByID(c.ID, c.cap); err != nil {
			t.Error(err)
		}

		// See if the entity can delete its target
		if err := em.DeleteEntity(c.ID, c.secret, c.toDelete); err != c.wantErr {
			t.Error(err)
		}
	}
}

func TestBasicCapabilities(t *testing.T) {
	em := New()

	s := []struct {
		ID         string
		uidNumber  int32
		secret     string
		capability string
		test       string
		err        error
	}{
		{"foo", -1, "a", "GLOBAL_ROOT", "GLOBAL_ROOT", nil},
		{"bar", -1, "a", "CREATE_ENTITY", "CREATE_ENTITY", nil},
		{"baz", -1, "a", "GLOBAL_ROOT", "CREATE_ENTITY", nil},
		{"bam", -1, "a", "CREATE_ENTITY", "GLOBAL_ROOT", E_ENTITY_UNQUALIFIED},
	}

	for _, c := range s {
		if err := em.newEntity(c.ID, c.uidNumber, c.secret); err != nil {
			t.Error(err)
		}

		e, err := em.getEntityByID(c.ID)
		if err != nil {
			t.Error(err)
		}

		em.setEntityCapability(e, c.capability)

		if err = em.checkEntityCapability(e, c.test); err != c.err {
			t.Error(err)
		}
	}
}

func TestSetSameCapabilityTwice(t *testing.T) {
	em := New()

	// Add an entity
	if err := em.newEntity("foo", -1, ""); err != nil {
		t.Error(err)
	}

	e, err := em.getEntityByID("foo")
	if err != nil {
		t.Error(err)
	}

	// Set one capability
	em.setEntityCapability(e, "GLOBAL_ROOT")
	if len(e.Meta.Capabilities) != 1 {
		t.Error("Wrong number of capabilities set!")
	}

	// Set it again and make sure its still only listed once.
	em.setEntityCapability(e, "GLOBAL_ROOT")
	if len(e.Meta.Capabilities) != 1 {
		t.Error("Wrong number of capabilities set!")
	}
}

func TestBasicCapabilitiesByID(t *testing.T) {
	em := New()

	s := []struct {
		ID         string
		uidNumber  int32
		secret     string
		capability string
		test       string
		err        error
	}{
		{"foo", -1, "a", "GLOBAL_ROOT", "GLOBAL_ROOT", nil},
		{"bar", -1, "a", "CREATE_ENTITY", "CREATE_ENTITY", nil},
		{"baz", -1, "a", "GLOBAL_ROOT", "CREATE_ENTITY", nil},
		{"bam", -1, "a", "CREATE_ENTITY", "GLOBAL_ROOT", E_ENTITY_UNQUALIFIED},
	}

	for _, c := range s {
		if err := em.newEntity(c.ID, c.uidNumber, c.secret); err != nil {
			t.Error(err)
		}

		em.setEntityCapabilityByID(c.ID, c.capability)

		if err := em.checkEntityCapabilityByID(c.ID, c.test); err != c.err {
			t.Error(err)
		}
	}
}

func TestCapabilityBogusEntity(t *testing.T) {
	em := New()

	// This test tries to set a capability on an entity that does
	// not exist.  In this case the error from getEntityByID
	// should be returned.
	if err := em.setEntityCapabilityByID("foo", "GLOBAL_ROOT"); err != E_NO_ENTITY {
		t.Error(err)
	}
}

func TestSetEntitySecretByID(t *testing.T) {
	em := New()

	s := []struct {
		ID        string
		uidNumber int32
		secret    string
	}{
		{"foo", 1, "a"},
		{"bar", 2, "a"},
		{"baz", 3, "a"},
	}

	// Load in the entities
	for _, c := range s {
		if err := em.newEntity(c.ID, c.uidNumber, c.secret); err != nil {
			t.Error(err)
		}
	}

	// Validate the secrets
	for _, c := range s {
		if err := em.ValidateSecret(c.ID, c.secret); err != nil {
			t.Errorf("Failed: want 'nil', got %v", err)
		}
	}
}

func TestSetEntitySecretByIDBogusEntity(t *testing.T) {
	em := New()

	// Attempt to set the secret on an entity that doesn't exist.
	if err := em.setEntitySecretByID("a", "a"); err != E_NO_ENTITY {
		t.Error(err)
	}
}

func TestValidateSecretBogusEntity(t *testing.T) {
	em := New()

	// Attempt to validate the secret on an entity that doesn't
	// exist, ensure that the right error is returned.
	if err := em.ValidateSecret("a", "a"); err != E_NO_ENTITY {
		t.Error(err)
	}
}

func TestValidateEntityCapabilityAndSecret(t *testing.T) {
	em := New()

	s := []struct {
		ID         string
		uidNumber  int32
		secret     string
		cap        string
		wantErr    error
		testSecret string
		testCap    string
	}{
		{"a", -1, "a", "", E_ENTITY_UNQUALIFIED, "a", "GLOBAL_ROOT"},
		{"b", -1, "a", "", E_ENTITY_BADAUTH, "b", ""},
		{"c", -1, "a", "CREATE_ENTITY", nil, "a", "CREATE_ENTITY"},
		{"d", -1, "a", "GLOBAL_ROOT", nil, "a", "CREATE_ENTITY"},
	}

	for _, c := range s {
		if err := em.newEntity(c.ID, c.uidNumber, c.secret); err != nil {
			t.Error(err)
		}

		if err := em.setEntityCapabilityByID(c.ID, c.cap); err != nil {
			t.Error(err)
		}

		if err := em.validateEntityCapabilityAndSecret(c.ID, c.testSecret, c.testCap); err != c.wantErr {
			t.Error(err)
		}
	}
}

func TestChangeSecret(t *testing.T) {
	em := New()

	entities := []struct {
		ID     string
		secret string
		cap    string
	}{
		{"a", "a", ""},                     // unpriviledged user
		{"b", "b", "CHANGE_ENTITY_SECRET"}, // can change other secrets
		{"c", "c", "GLOBAL_ROOT"},          // can also change other secrets
	}

	cases := []struct {
		ID           string
		secret       string
		changeID     string
		changeSecret string
		wantErr      error
	}{
		{"a", "e", "a", "a", E_ENTITY_BADAUTH},     // same entity, bad secret
		{"a", "a", "a", "b", nil},                  // can change own password
		{"a", "b", "b", "d", E_ENTITY_UNQUALIFIED}, // can't change other secrets without capability
		{"b", "b", "a", "a", nil},                  // can change other's secret with CHANGE_ENTITY_SECRET
		{"c", "c", "a", "b", nil},                  // can change other's secret with GLOBAL_ROOT
	}

	// Add some entities
	for _, e := range entities {
		if err := em.newEntity(e.ID, -1, e.secret); err != nil {
			t.Error(err)
		}

		if err := em.setEntityCapabilityByID(e.ID, e.cap); err != nil {
			t.Error(err)
		}
	}

	// Run the tests
	for _, c := range cases {
		if err := em.ChangeSecret(c.ID, c.secret, c.changeID, c.changeSecret); err != c.wantErr {
			t.Error(err)
		}
	}
}
