package entity_manager

import "testing"

func TestBasicCapabilities(t *testing.T) {
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

	resetMap()

	for _, c := range s {
		if err := newEntity(c.ID, c.uidNumber, c.secret); err != nil {
			t.Error(err)
		}

		e, err := getEntityByID(c.ID)
		if err != nil {
			t.Error(err)
		}

		setEntityCapability(e, c.capability)

		if err = checkEntityCapability(e, c.test); err != c.err {
			t.Error(err)
		}
	}
}

func TestSetSameCapabilityTwice(t *testing.T) {
	// Reset state
	resetMap()

	// Add an entity
	if err := newEntity("foo", -1, ""); err != nil {
		t.Error(err)
	}

	e, err := getEntityByID("foo")
	if err != nil {
		t.Error(err)
	}

	// Set one capability
	setEntityCapability(e, "GLOBAL_ROOT")
	if len(e.Meta.Capabilities) != 1 {
		t.Error("Wrong number of capabilities set!")
	}

	// Set it again and make sure its still only listed once.
	setEntityCapability(e, "GLOBAL_ROOT")
	if len(e.Meta.Capabilities) != 1 {
		t.Error("Wrong number of capabilities set!")
	}
}

func TestBasicCapabilitiesByID(t *testing.T) {
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

	resetMap()

	for _, c := range s {
		if err := newEntity(c.ID, c.uidNumber, c.secret); err != nil {
			t.Error(err)
		}

		setEntityCapabilityByID(c.ID, c.capability)

		if err := checkEntityCapabilityByID(c.ID, c.test); err != c.err {
			t.Error(err)
		}
	}
}

func TestCapabilityBogusEntity(t *testing.T) {
	// This test tries to set a capability on an entity that does
	// not exist.  In this case the error from getEntityByID
	// should be returned.
	resetMap()
	if err := setEntityCapabilityByID("foo", "GLOBAL_ROOT"); err != E_NO_ENTITY {
		t.Error(err)
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
		if err := newEntity(c.ID, c.uidNumber, c.secret); err != nil {
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

func TestSetEntitySecretByIDBogusEntity(t *testing.T) {
	// Attempt to set the secret on an entity that doesn't exist.
	resetMap()
	if err := setEntitySecretByID("a", "a"); err != E_NO_ENTITY {
		t.Error(err)
	}
}

func TestValidateEntitySecretByIDBogusEntity(t *testing.T) {
	// Attempt to validate the secret on an entity that doesn't
	// exist, ensure that the right error is returned.
	resetMap()
	if err := ValidateEntitySecretByID("a", "a"); err != E_NO_ENTITY {
		t.Error(err)
	}
}

func TestValidateEntityCapabilityAndSecret(t *testing.T) {
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

	resetMap()

	for _, c := range s {
		if err := newEntity(c.ID, c.uidNumber, c.secret); err != nil {
			t.Error(err)
		}

		if err := setEntityCapabilityByID(c.ID, c.cap); err != nil {
			t.Error(err)
		}

		if err := validateEntityCapabilityAndSecret(c.ID, c.testSecret, c.testCap); err != c.wantErr {
			t.Error(err)
		}
	}
}

func TestChangeSecret(t *testing.T) {
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
		{"a", "e", "a", "a", E_ENTITY_BADAUTH},           // same entity, bad secret
		{"a", "a", "a", "b", nil},                  // can change own password
		{"a", "b", "b", "d", E_ENTITY_UNQUALIFIED}, // can't change other secrets without capability
		{"b", "b", "a", "a", nil},                  // can change other's secret with CHANGE_ENTITY_SECRET
		{"c", "c", "a", "b", nil},                  // can change other's secret with GLOBAL_ROOT
	}

	resetMap()

	// Add some entities
	for _, e := range entities {
		if err := newEntity(e.ID, -1, e.secret); err != nil {
			t.Error(err)
		}

		if err := setEntityCapabilityByID(e.ID, e.cap); err != nil {
			t.Error(err)
		}
	}

	// Run the tests
	for _, c := range cases {
		if err := ChangeSecret(c.ID, c.secret, c.changeID, c.changeSecret); err != c.wantErr {
			t.Error(err)
		}
	}
}
