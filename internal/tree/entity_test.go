package tree

import (
	"testing"

	"github.com/NetAuth/NetAuth/internal/db/impl/MemDB"
	"github.com/NetAuth/NetAuth/internal/crypto/impl/nocrypto"
	"github.com/NetAuth/NetAuth/pkg/errors"
	"github.com/golang/protobuf/proto"

	pb "github.com/NetAuth/Protocol"
)

func TestAddDuplicateID(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

	s := []struct {
		ID        string
		number int32
		secret    string
		err       error
	}{
		{"foo", 1, "", nil},
		{"foo", 2, "", errors.E_DUPLICATE_ID},
	}

	for _, c := range s {
		if err := em.NewEntity(c.ID, c.number, c.secret); err != c.err {
			t.Errorf("Got %v; Want: %v", err, c.err)
		}
	}
}

func TestAddDuplicateUIDNumber(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

	s := []struct {
		ID        string
		number int32
		secret    string
		err       error
	}{
		{"foo", 1, "", nil},
		{"bar", 1, "", errors.E_DUPLICATE_UIDNUMBER},
	}

	for _, c := range s {
		if err := em.NewEntity(c.ID, c.number, c.secret); err != c.err {
			t.Errorf("Got %v; Want: %v", err, c.err)
		}
	}
}

func TestNewEntityAutoNumber(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

	s := []struct {
		ID        string
		number int32
		secret    string
	}{
		{"foo", 1, ""},
		{"bar", -1, ""},
		{"baz", 3, ""},
	}

	for _, c := range s {
		if err := em.NewEntity(c.ID, c.number, c.secret); err != nil {
			t.Error(err)
		}
	}
}

func TestDeleteEntityByID(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

	s := []struct {
		ID        string
		number int32
		secret    string
	}{
		{"foo", 1, ""},
		{"bar", 2, ""},
		{"baz", 3, ""},
	}

	// Populate some entities
	for _, c := range s {
		if err := em.NewEntity(c.ID, c.number, c.secret); err != nil {
			t.Error(err)
		}
	}

	for _, c := range s {
		// Delete the entity
		if err := em.DeleteEntityByID(c.ID); err != nil {
			t.Error(err)
		}

		// Make sure checking for that entity returns errors.E_NO_ENTITY
		if _, err := em.db.LoadEntity(c.ID); err != errors.E_NO_ENTITY {
			t.Error(err)
		}
	}
}

func TestSetSameCapabilityTwice(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

	// Add an entity
	if err := em.NewEntity("foo", -1, ""); err != nil {
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

func TestCapabilityBogusEntity(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

	// This test tries to set a capability on an entity that does
	// not exist.  In this case the error from getEntityByID
	// should be returned.
	if err := em.SetEntityCapabilityByID("foo", "GLOBAL_ROOT"); err != errors.E_NO_ENTITY {
		t.Error(err)
	}
}

func TestSetEntitySecretByID(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

	s := []struct {
		ID        string
		number int32
		secret    string
	}{
		{"foo", 1, "a"},
		{"bar", 2, "a"},
		{"baz", 3, "a"},
	}

	// Load in the entities
	for _, c := range s {
		if err := em.NewEntity(c.ID, c.number, c.secret); err != nil {
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
	if err := em.SetEntitySecretByID("a", "a"); err != errors.E_NO_ENTITY {
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

func TestGetEntity(t *testing.T) {
	em := New(MemDB.New(), nocrypto.New())

	// Add a new entity with known values.
	if err := em.NewEntity("foo", -1, ""); err != nil {
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
		Number: proto.Int32(1),
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
	if err := em.NewEntity("foo", -1, ""); err != nil {
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
	em.UpdateEntityMeta(e.GetID(), fullMeta)

	// Verify that the update above took
	if e.GetMeta().GetLegalName() != "Foobert McMillan" {
		t.Error("Field set mismatch!")
	}

	// This is metadata that can't be updated with this call,
	// verify that it gets dropped.
	groups := []string{"fooGroup"}
	badMeta := &pb.EntityMeta{
		Groups: groups,
	}
	em.UpdateEntityMeta(e.GetID(), badMeta)

	// The update from badMeta should not have gone through, and
	// the old value should still be present.
	if e.GetMeta().Groups != nil {
		t.Errorf("badMeta was merged! (%v)", e.GetMeta().GetGroups())
	}
	if e.GetMeta().GetLegalName() != "Foobert McMillan" {
		t.Error("Update overwrote unset value!")
	}
}
