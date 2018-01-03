package entity_manager

import (
	"log"

	pb "github.com/NetAuth/NetAuth/proto"
)

var (
	// eByID is a package scoped map of ID strings to entities.
	eByID = make(map[string]*pb.Entity)

	// eByUIDNumber is a package scoped map of int32 entity
	// uidNumbers to entities.
	eByUIDNumber = make(map[int32]*pb.Entity)

	// Making a bootstrap entity is a rare thing and short
	// circuits most of the permissions logic.  As such we only
	// allow it to be done once per server start.
	bootstrap_done bool = false
)

// nextUIDNumber computes the next available uidNumber to be assigned.
// This allows a NewEntity request to be made with the uidNumber field
// unset.
func nextUIDNumber() int32 {
	var largest int32 = 0

	// Iterate over the entities and return the largest ID found
	// +1.  This allows them to be in any order or have IDs
	// missing in the middle and still work.  Though an
	// inefficient search this is worst case O(N) and happens in
	// memory anyway.
	for i := range eByID {
		if eByID[i].GetUidNumber() > largest {
			largest = eByID[i].GetUidNumber()
		}
	}

	return largest + 1
}

// newEntity creates a new entity given an ID, uidNumber, and secret.
// Its not necessary to set the secret upon creation and it can be set
// later.  If not set on creaction then the entity will not be usable.
// uidNumber must be a unique positive integer.  Because these are
// generally allocated in sequence the special value '-1' may be
// specified which will select the next available number.
func newEntity(ID string, uidNumber int32, secret string) error {
	// Does this entity exist already?
	_, exists := eByID[ID]
	if exists {
		log.Printf("Entity with ID '%s' already exists!", ID)
		return E_DUPLICATE_ID
	}
	_, exists = eByUIDNumber[uidNumber]
	if exists {
		log.Printf("Entity with uidNumber '%d' already exists!", uidNumber)
		return E_DUPLICATE_UIDNUMBER
	}

	// Were we given a specific uidNumber?
	if uidNumber == -1 {
		// -1 is a sentinel value that tells us to pick the
		// next available number and assign it.
		uidNumber = nextUIDNumber()
	}

	// Ok, they don't exist so we'll make them exist now
	newEntity := &pb.Entity{
		ID:        &ID,
		UidNumber: &uidNumber,
		Secret:    &secret,
		Meta:      &pb.EntityMeta{},
	}

	// Add this entity to the in-memory listings
	eByID[ID] = newEntity
	eByUIDNumber[uidNumber] = newEntity

	// Successfully created we now return no errors
	log.Printf("Created entity '%s'", ID)

	// Now we set the entity secret, this could be inlined, but
	// having it in the seperate function makes resetting the
	// secret trivial.
	setEntitySecretByID(ID, secret)

	return nil
}

// NewEntity is a public function which adds a new entity on behalf of
// another one.  The requesting entity must be able to validate its
// identity and posses the appropriate capability to add a new entity
// to the system.
func NewEntity(requestID, requestSecret, newID string, newUIDNumber int32, newSecret string) error {
	// Validate that the entity is real and permitted to perform
	// this action.
	if err := validateEntityCapabilityAndSecret(requestID, requestSecret, "CREATE_ENTITY"); err != nil {
		return err
	}

	// The entity is who they say they are and has the appropriate
	// capability, time to actually create the new entity.
	if err := newEntity(newID, newUIDNumber, newSecret); err != nil {
		return err
	}
	return nil
}

// NewBootstrapEntity is a function that can be called during the
// startup of the srever to create an entity that has the appropriate
// authority to create more entities and otherwise manage the server.
// This can only be called once during startup, attepts to call it
// again will result in no change.  The bootstrap user will always get
// the next available number which in most cases will be 1.
func MakeBootstrap(ID string, secret string) {
	if bootstrap_done {
		return
	}

	// In some cases if there is an existing system that has no
	// admin, it is necessary to confer bootstrap powers to an
	// existing user.  In that case they are just selected and
	// then provided the GLOBAL_ROOT capability.
	e, err := getEntityByID(ID)
	if err != nil {
		log.Printf("No entity with ID '%s' exists!  Creating...", ID)
	}

	// This is not a normal Go way of doing this, but this
	// function has two possible success cases, the flow may jump
	// in here and return if there is an existing entity to get
	// root powers.
	if e != nil {
		setEntityCapability(e, "GLOBAL_ROOT")
		bootstrap_done = true
		return
	}

	// Even in the bootstrap case its still possible this can
	// fail, in that case its useful to have the error.
	if err := newEntity(ID, -1, secret); err != nil {
		log.Printf("Could not create bootstrap user! (%s)", err)
	}
	if err := setEntityCapabilityByID(ID, "GLOBAL_ROOT"); err != nil {
		log.Printf("Couldn't provide root authority! (%s)", err)
	}

	bootstrap_done = true
}

// DisableBootstrap disables the ability to bootstrap after the
// opportunity to do so has passed.
func DisableBootstrap() {
	bootstrap_done = true
}

// getEntityByID returns a pointer to an Entity struct and an error
// value.  The error value will either be E_NO_ENTITY if the requested
// value did not match, or will be nil where an entity is returned.
// The string must be a complete match for the entity name being
// requested.
func getEntityByID(ID string) (*pb.Entity, error) {
	e, ok := eByID[ID]
	if !ok {
		return nil, E_NO_ENTITY
	}
	return e, nil
}

// GetEntityByUIDNumber returns a pointer to an Entity struct and an
// error value.  The error value will either be E_NO_ENTITY if the
// requested value did not match, or will be nil where an entity is
// returned.  The numeric value must be an exact match for the entity
// being requested in addition to being an int32 value.
func getEntityByUIDNumber(uidNumber int32) (*pb.Entity, error) {
	e, ok := eByUIDNumber[uidNumber]
	if !ok {
		return nil, E_NO_ENTITY
	}
	return e, nil
}

// DeleteEntityByID deletes the named entity.  This function will
// delete the entity in a non-atomic way, but will ensure that the
// entity cannot be authenticated with before returning.  If the named
// ID does not exist the function will return E_NO_ENTITY, in all
// other cases nil is returned.
func deleteEntityByID(ID string) error {
	e, err := getEntityByID(ID)
	if err != nil {
		return E_NO_ENTITY
	}

	// There's a small chance that the deletes won't go through
	// but if that happens there's not much that we can do about
	// it since delete() is a language construct and there's no
	// way to drill down deeper.  To be paranoid though we'll
	// clear the secret on the entity which should prevent it from
	// being usable.
	e.Secret = nil

	// Now we need to delete from both maps
	delete(eByID, e.GetID())
	delete(eByUIDNumber, e.GetUidNumber())
	log.Printf("Deleted entity '%s'", e.GetID())

	return nil
}

func DeleteEntity(requestID string, requestSecret string, deleteID string) error {
	// Validate that the entity is real and permitted to perform
	// this action.
	if err := validateEntityCapabilityAndSecret(requestID, requestSecret, "DELETE_ENTITY"); err != nil {
		return err
	}

	// Delete the requested entity
	return deleteEntityByID(deleteID)
}
