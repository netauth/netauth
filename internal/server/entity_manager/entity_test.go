package entity_manager

import (
	"testing"

	"github.com/NetAuth/NetAuth/internal/server/db/impl/MemDB"
	"github.com/NetAuth/NetAuth/internal/server/crypto/impl/nocrypto"
	"github.com/NetAuth/NetAuth/pkg/errors"
	"github.com/golang/protobuf/proto"

	pb "github.com/NetAuth/NetAuth/pkg/proto"
)

func TestNewEntityWithAuth(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

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
		{"b", -1, "a", "", "bb", -1, "a", errors.E_ENTITY_UNQUALIFIED},
		{"c", -1, "a", "CREATE_ENTITY", "cc", -1, "a", nil},
		{"d", -1, "a", "CREATE_ENTITY", "a", -1, "a", errors.E_DUPLICATE_ID},
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
	em := New(MemDB.New(), nocrypto.New())

	s := []struct {
		ID        string
		uidNumber int32
		secret    string
		err       error
	}{
		{"foo", 1, "", nil},
		{"foo", 2, "", errors.E_DUPLICATE_ID},
	}

	for _, c := range s {
		if err := em.newEntity(c.ID, c.uidNumber, c.secret); err != c.err {
			t.Errorf("Got %v; Want: %v", err, c.err)
		}
	}
}

func TestAddDuplicateUIDNumber(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

	s := []struct {
		ID        string
		uidNumber int32
		secret    string
		err       error
	}{
		{"foo", 1, "", nil},
		{"bar", 1, "", errors.E_DUPLICATE_UIDNUMBER},
	}

	for _, c := range s {
		if err := em.newEntity(c.ID, c.uidNumber, c.secret); err != c.err {
			t.Errorf("Got %v; Want: %v", err, c.err)
		}
	}
}

func TestNewEntityAutoNumber(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

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
	em := New(MemDB.New(), nocrypto.New())

	s := []struct {
		ID            string
		secret        string
		pre_disable   bool
		bootstrap_val bool
		wantErr       error
	}{
		{"bar", "", false, false, nil},
		{"baz", "", false, false, nil},
		{"quu", "", true, true, errors.E_NO_ENTITY},
		{"qux", "", true, false, errors.E_NO_ENTITY},
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

func TestDeleteEntityByID(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

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

		// Make sure checking for that entity returns errors.E_NO_ENTITY
		if _, err := em.db.LoadEntity(c.ID); err != errors.E_NO_ENTITY {
			t.Error(err)
		}
	}
}

func TestDeleteEntityWithAuth(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

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
		{"b", -1, "a", "", "bar", errors.E_ENTITY_UNQUALIFIED},
		{"c", -1, "a", "CREATE_ENTITY", "baz", errors.E_ENTITY_UNQUALIFIED},
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
	em := New(MemDB.New(), nocrypto.New())

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
		{"bam", -1, "a", "CREATE_ENTITY", "GLOBAL_ROOT", errors.E_ENTITY_UNQUALIFIED},
	}

	for _, c := range s {
		if err := em.newEntity(c.ID, c.uidNumber, c.secret); err != nil {
			t.Error(err)
		}

		e, err := em.db.LoadEntity(c.ID)
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
	em := New(MemDB.New(), nocrypto.New())

	// Add an entity
	if err := em.newEntity("foo", -1, ""); err != nil {
		t.Error(err)
	}

	e, err := em.db.LoadEntity("foo")
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
	em := New(MemDB.New(), nocrypto.New())

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
		{"bam", -1, "a", "CREATE_ENTITY", "GLOBAL_ROOT", errors.E_ENTITY_UNQUALIFIED},
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
	em := New(MemDB.New(), nocrypto.New())

	// This test tries to set a capability on an entity that does
	// not exist.  In this case the error from getEntityByID
	// should be returned.
	if err := em.setEntityCapabilityByID("foo", "GLOBAL_ROOT"); err != errors.E_NO_ENTITY {
		t.Error(err)
	}
}

func TestSetEntitySecretByID(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

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
	em := New(MemDB.New(), nocrypto.New())

	// Attempt to set the secret on an entity that doesn't exist.
	if err := em.setEntitySecretByID("a", "a"); err != errors.E_NO_ENTITY {
		t.Error(err)
	}
}

func TestValidateSecretBogusEntity(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

	// Attempt to validate the secret on an entity that doesn't
	// exist, ensure that the right error is returned.
	if err := em.ValidateSecret("a", "a"); err != errors.E_NO_ENTITY {
		t.Error(err)
	}
}

func TestValidateEntityCapabilityAndSecret(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

	s := []struct {
		ID         string
		uidNumber  int32
		secret     string
		cap        string
		wantErr    error
		testSecret string
		testCap    string
	}{
		{"a", -1, "a", "", errors.E_ENTITY_UNQUALIFIED, "a", "GLOBAL_ROOT"},
		{"b", -1, "a", "", errors.E_ENTITY_BADAUTH, "b", ""},
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
	em := New(MemDB.New(), nocrypto.New())

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
		{"a", "e", "a", "a", errors.E_ENTITY_BADAUTH},     // same entity, bad secret
		{"a", "a", "a", "b", nil},                         // can change own password
		{"a", "b", "b", "d", errors.E_ENTITY_UNQUALIFIED}, // can't change other secrets without capability
		{"b", "b", "a", "a", nil},                         // can change other's secret with CHANGE_ENTITY_SECRET
		{"c", "c", "a", "b", nil},                         // can change other's secret with GLOBAL_ROOT
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

func TestGetEntity(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

	// Add a new entity with known values.
	if err := em.newEntity("foo", -1, ""); err != nil {
		t.Error(err)
	}

	// First validate that this works with no entity
	entity, err := em.GetEntity("")
	if err != errors.E_NO_ENTITY {
		t.Error(err)
	}

	// Now check that we get back the right values for the entity
	// we added earlier.
	entity, err = em.GetEntity("foo")
	if err != nil {
		t.Error(err)
	}

	entityTest := &pb.Entity{
		ID:        proto.String("foo"),
		UidNumber: proto.Int32(1),
		Secret:    proto.String("<REDACTED>"),
		Meta:      &pb.EntityMeta{},
	}

	if !proto.Equal(entity, entityTest) {
		t.Errorf("Entity retrieved not equal! got %v want %v", entity, entityTest)
	}
}

func TestUpdateEntityMetaInternal(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

	// Add a new entity with known values
	if err := em.newEntity("foo", -1, ""); err != nil {
		t.Error(err)
	}

	fullMeta := &pb.EntityMeta{
		LegalName: proto.String("Foobert McMillan"),
	}

	// This checks that merging into the empty default meta works,
	// since these will be the only values set.
	e, err := em.db.LoadEntity("foo")
	if err != nil {
		t.Error(err)
	}
	em.updateEntityMeta(e, fullMeta)

	// Verify that the update above took
	if e.GetMeta().GetLegalName() != "Foobert McMillan" {
		t.Error("Field set mismatch!")
	}

	// This is metadata that can't be updated with this call,
	// verify that it gets dropped.
	groups := []*pb.Group{&pb.Group{}}
	badMeta := &pb.EntityMeta{
		Groups: groups,
	}
	em.updateEntityMeta(e, badMeta)

	// The update from badMeta should not have gone through, and
	// the old value should still be present.
	if e.GetMeta().Groups != nil {
		t.Errorf("badMeta was merged! (%v)", e.GetMeta().GetGroups())
	}
	if e.GetMeta().GetLegalName() != "Foobert McMillan" {
		t.Error("Update overwrote unset value!")
	}
}

func TestUpdateEntityMetaExternal(t *testing.T) {
	s := []struct {
		ID         string
		secret     string
		capability string
		modID      string
		wantErr    error
	}{
		{"foo", "foo", "", "a", errors.E_ENTITY_UNQUALIFIED},
		{"foo", "", "", "a", errors.E_ENTITY_BADAUTH},
		{"foo", "foo", "MODIFY_ENTITY_META", "a", nil},
		{"a", "b", "", "a", errors.E_ENTITY_BADAUTH},
		{"a", "a", "", "a", nil},
	}

	em := New(MemDB.New(), nocrypto.New())

	if err := em.newEntity("foo", -1, "foo"); err != nil {
		t.Error(err)
	}
	if err := em.newEntity("a", -1, "a"); err != nil {
		t.Error(err)
	}

	modMeta := &pb.EntityMeta{DisplayName: proto.String("Waldo")}

	for _, c := range s {
		if err := em.setEntityCapabilityByID(c.ID, c.capability); err != nil {
			t.Error(err)
		}
		if err := em.UpdateEntityMeta(c.ID, c.secret, c.modID, modMeta); err != c.wantErr {
			t.Error(err)
		}
	}
}
