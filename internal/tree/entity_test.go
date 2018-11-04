package tree

import (
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/NetAuth/NetAuth/internal/crypto"
	"github.com/NetAuth/NetAuth/internal/db"

	pb "github.com/NetAuth/Protocol"
)

func TestNextUIDNumber(t *testing.T) {
	em := getNewEntityManager(t)

	s := []struct {
		ID            string
		number        int32
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

func TestAddDuplicateID(t *testing.T) {
	em := getNewEntityManager(t)

	s := []struct {
		ID     string
		number int32
		secret string
		err    error
	}{
		{"foo", 1, "", nil},
		{"foo", 2, "", ErrDuplicateEntityID},
	}

	for _, c := range s {
		if err := em.NewEntity(c.ID, c.number, c.secret); err != c.err {
			t.Errorf("Got %v; Want: %v", err, c.err)
		}
	}
}

func TestNewEntityAutoNumber(t *testing.T) {
	em := getNewEntityManager(t)

	s := []struct {
		ID     string
		number int32
		secret string
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

func TestMakeBootstrapDoubleBootstrap(t *testing.T) {
	em := getNewEntityManager(t)

	// Claim the bootstrap is already done
	em.bootstrapDone = true
	em.MakeBootstrap("", "")
}

func TestMakeBootstrapExtantEntity(t *testing.T) {
	em := getNewEntityManager(t)

	if err := em.NewEntity("foo", -1, "foo"); err != nil {
		t.Fatal(err)
	}

	em.MakeBootstrap("foo", "foo")

	e, err := em.GetEntity("foo")
	if err != nil {
		t.Fatal(err)
	}

	gRoot := pb.Capability(pb.Capability_value["GLOBAL_ROOT"])

	if e.GetMeta().GetCapabilities()[0] != gRoot {
		t.Fatalf("Unexpected capability: %s", e.GetMeta().GetCapabilities())
	}
}

func TestMakeBootstrapCreateEntity(t *testing.T) {
	em := getNewEntityManager(t)

	em.MakeBootstrap("foo", "foo")

	e, err := em.GetEntity("foo")
	if err != nil {
		t.Fatal(err)
	}

	gRoot := pb.Capability(pb.Capability_value["GLOBAL_ROOT"])

	if e.GetMeta().GetCapabilities()[0] != gRoot {
		t.Fatalf("Unexpected capability: %s", e.GetMeta().GetCapabilities())
	}
}

func TestBootstrapLockedEntity(t *testing.T) {
	em := getNewEntityManager(t)

	if err := em.NewEntity("foo", -1, "foo"); err != nil {
		t.Fatal(err)
	}

	if err := em.LockEntity("foo"); err != nil {
		t.Fatal(err)
	}

	if err := em.ValidateSecret("foo", "foo"); err != ErrEntityLocked {
		t.Fatal(err)
	}

	em.MakeBootstrap("foo", "foo")

	if err := em.ValidateSecret("foo", "foo"); err != nil {
		t.Fatal(err)
	}
}

func TestDisableBootstrap(t *testing.T) {
	em := getNewEntityManager(t)

	if em.bootstrapDone == true {
		t.Fatal("Bootstrap is somehow already done")
	}
	em.DisableBootstrap()
	if em.bootstrapDone == false {
		t.Fatal("Bootstrap somehow not done")
	}
}

func TestDeleteEntityByID(t *testing.T) {
	em := getNewEntityManager(t)

	s := []struct {
		ID     string
		number int32
		secret string
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

		// Make sure checking for that entity returns db.ErrUnknownEntity
		if _, err := em.db.LoadEntity(c.ID); err != db.ErrUnknownEntity {
			t.Error(err)
		}
	}
}

func TestDeleteEntityAgain(t *testing.T) {
	em := getNewEntityManager(t)
	if err := em.DeleteEntityByID("foo"); err != db.ErrUnknownEntity {
		t.Fatalf("Wrong error: %s", err)
	}
}

func TestSetSameCapabilityTwice(t *testing.T) {
	em := getNewEntityManager(t)

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

func TestSetCapabilityBogusEntity(t *testing.T) {
	em := getNewEntityManager(t)

	// This test tries to set a capability on an entity that does
	// not exist.  In this case the error from getEntityByID
	// should be returned.
	if err := em.SetEntityCapabilityByID("foo", "GLOBAL_ROOT"); err != db.ErrUnknownEntity {
		t.Error(err)
	}
}

func TestSetCapabilityNoCap(t *testing.T) {
	em := getNewEntityManager(t)

	if err := em.NewEntity("foo", -1, "foo"); err != nil {
		t.Fatal(err)
	}

	if err := em.SetEntityCapabilityByID("foo", ""); err != ErrUnknownCapability {
		t.Error(err)
	}
}

func TestRemoveCapability(t *testing.T) {
	em := getNewEntityManager(t)

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
	// Set another capability
	em.setEntityCapability(e, "MODIFY_ENTITY_META")
	if len(e.Meta.Capabilities) != 2 {
		t.Error("Wrong number of capabilities set!")
	}

	// Remove it and make sure its gone
	em.removeEntityCapability(e, "GLOBAL_ROOT")
	if len(e.Meta.Capabilities) != 1 {
		t.Error("Wrong number of capabilities set!")
	}
}

func TestRemoveCapabilityBogusEntity(t *testing.T) {
	em := getNewEntityManager(t)

	if err := em.RemoveEntityCapabilityByID("foo", "GLOBAL_ROOT"); err != db.ErrUnknownEntity {
		t.Error(err)
	}
}

func TestRemoveCapabilityNoCap(t *testing.T) {
	em := getNewEntityManager(t)

	if err := em.NewEntity("foo", -1, "foo"); err != nil {
		t.Fatal(err)
	}

	if err := em.RemoveEntityCapabilityByID("foo", ""); err != ErrUnknownCapability {
		t.Error(err)
	}
}

func TestSetEntitySecretByID(t *testing.T) {
	em := getNewEntityManager(t)

	s := []struct {
		ID     string
		number int32
		secret string
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
	em := getNewEntityManager(t)

	// Attempt to set the secret on an entity that doesn't exist.
	if err := em.SetEntitySecretByID("a", "a"); err != db.ErrUnknownEntity {
		t.Error(err)
	}
}

func TestValidateSecretBogusEntity(t *testing.T) {
	em := getNewEntityManager(t)

	// Attempt to validate the secret on an entity that doesn't
	// exist, ensure that the right error is returned.
	if err := em.ValidateSecret("a", "a"); err != db.ErrUnknownEntity {
		t.Error(err)
	}
}

func TestValidateSecretWrongSecret(t *testing.T) {
	em := getNewEntityManager(t)

	if err := em.NewEntity("foo", -1, "foo"); err != nil {
		t.Fatal(err)
	}

	if err := em.ValidateSecret("foo", "bar"); err != crypto.ErrAuthorizationFailure {
		t.Fatal(err)
	}
}

func TestGetEntity(t *testing.T) {
	em := getNewEntityManager(t)

	// Add a new entity with known values.
	if err := em.NewEntity("foo", -1, ""); err != nil {
		t.Error(err)
	}

	// First validate that this works with no entity
	_, err := em.GetEntity("")
	if err != db.ErrUnknownEntity {
		t.Error(err)
	}

	// Now check that we get back the right values for the entity
	// we added earlier.
	entity, err := em.GetEntity("foo")
	if err != nil {
		t.Error(err)
	}

	entityTest := &pb.Entity{
		ID:     proto.String("foo"),
		Number: proto.Int32(1),
		Secret: proto.String("<REDACTED>"),
	}

	if !proto.Equal(entity, entityTest) {
		t.Errorf("Entity retrieved not equal! got %v want %v", entity, entityTest)
	}
}

func TestUpdateEntityMetaInternal(t *testing.T) {
	em := getNewEntityManager(t)

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

func TestUpdateEntityMetaExternalNoEntity(t *testing.T) {
	em := getNewEntityManager(t)

	if err := em.UpdateEntityMeta("non-existent", nil); err != db.ErrUnknownEntity {
		t.Fatal(err)
	}
}

func TestUpdateEntityKeys(t *testing.T) {
	em := getNewEntityManager(t)

	if err := em.NewEntity("foo", -1, "bar"); err != nil {
		t.Error(err)
	}

	if _, err := em.UpdateEntityKeys("foo", "ADD", "SIMPLE", "KEYCODE"); err != nil {
		t.Error(err)
	}

	l, err := em.UpdateEntityKeys("foo", "LIST", "", "")
	if err != nil {
		t.Error(err)
	}
	if len(l) != 1 || l[0] != "SIMPLE:KEYCODE" {
		t.Errorf("Bad Key: %v", l)
	}

	if _, err := em.UpdateEntityKeys("foo", "DEL", "", "KEY"); err != nil {
		t.Error(err)
	}

	l, err = em.UpdateEntityKeys("foo", "LIST", "", "")
	if err != nil {
		t.Error(err)
	}
	if len(l) != 0 {
		t.Errorf("Zombie keys: %s", l)
	}
}

func TestManageUntypedEntityMeta(t *testing.T) {
	em := getNewEntityManager(t)

	if err := em.NewEntity("foo", -1, ""); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		entityID string
		mode     string
		key      string
		value    string
		wantKV   []string
		wantErr  error
	}{
		{"foo", "upsert", "k1", "v1", nil, nil},
		{"foo", "read", "*", "", []string{"k1:v1"}, nil},
		{"unknown", "read", "*", "", nil, db.ErrUnknownEntity},
	}

	for i, c := range cases {
		kv, err := em.ManageUntypedEntityMeta(c.entityID, c.mode, c.key, c.value)
		if err != c.wantErr {
			t.Fatalf("%d: Got: %v; Want: %v", i, err, c.wantErr)
		}

		// This function is defined in helpers_test.go
		if !slicesAreEqual(kv, c.wantKV) {
			t.Fatalf("%d: Got: %v; Want: %v", i, kv, c.wantKV)
		}
	}
}

func TestLockUnlockEntity(t *testing.T) {
	em := getNewEntityManager(t)
	if err := em.NewEntity("foo", -1, "bar"); err != nil {
		t.Fatal(err)
	}

	if err := em.LockEntity("foo"); err != nil {
		t.Fatal(err)
	}

	if err := em.ValidateSecret("foo", "bar"); err != ErrEntityLocked {
		t.Fatal(err)
	}

	if err := em.UnlockEntity("foo"); err != nil {
		t.Fatal(err)
	}

	if err := em.ValidateSecret("foo", "bar"); err != nil {
		t.Fatal(err)
	}
}

func TestSetEntityLockState(t *testing.T) {
	em := getNewEntityManager(t)

	if err := em.setEntityLockState("does-not-exist", true); err != db.ErrUnknownEntity {
		t.Fatal(err)
	}
}

func TestSafeCopyEntity(t *testing.T) {
	em := getNewEntityManager(t)

	if err := em.NewEntity("foo", -1, "bar"); err != nil {
		t.Error(err)
	}

	e, err := em.db.LoadEntity("foo")
	if err != nil {
		t.Error(err)
	}

	ne := safeCopyEntity(e)

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
