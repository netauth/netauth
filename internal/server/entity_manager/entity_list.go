package entity_manager

import (
	"log"

	"golang.org/x/crypto/bcrypt"

	pb "github.com/NetAuth/NetAuth/proto"
)

var (
	// eByID is a package scoped map of ID strings to entities.
	eByID = make(map[string]*pb.Entity)

	// eByUIDNumber is a package scoped map of int32 entity
	// uidNumbers to entities.
	eByUIDNumber = make(map[int32]*pb.Entity)
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

// NewEntity creates a new entity given an ID, uidNumber, and secret.
// Its not necessary to set the secret upon creation and it can be set
// later.  If not set on creaction then the entity will not be usable.
// uidNumber must be a unique positive integer.  Because these are
// generally allocated in sequence the special value '-1' may be
// specified which will select the next available number.
func NewEntity(ID string, uidNumber int32, secret string) error {
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
	SetEntitySecretByID(ID, secret)

	return nil
}

// GetEntityByID returns a pointer to an Entity struct and an error
// value.  The error value will either be E_NO_ENTITY if the requested
// value did not match, or will be nil where an entity is returned.
// The string must be a complete match for the entity name being
// requested.
func GetEntityByID(ID string) (*pb.Entity, error) {
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
func GetEntityByUIDNumber(uidNumber int32) (*pb.Entity, error) {
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
func DeleteEntityByID(ID string) error {
	e, err := GetEntityByID(ID)
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

// DeleteEntityByUIDNumber deletes the numbered entity.  This function
// will delete the entity in a non-atomic way, but will ensure that
// the entity cannot be authenticated with before returning.  If the
// provided uidNumber does not exist the function will return
// E_NO_ENTITY, in all other cases nil is returned.  The uidNumber
// must be of type int32.
func DeleteEntityByUIDNumber(n int32) error {
	e, err := GetEntityByUIDNumber(n)
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

// SetEntitySecretByID sets the secret on a given entity using the
// bcrypt secure hashing algorithm.
func SetEntitySecretByID(ID string, secret string) error {
	e, err := GetEntityByID(ID)
	if err != nil {
		return err
	}

	// todo(maldridge) This is currently configured to use the
	// minimum cost hash.  THIS IS NOT SECURE.  Before releasing
	// 1.0 this must be changes to pull this value either
	// dynamically or from a config file somewhere.
	hash, err := bcrypt.GenerateFromPassword([]byte(secret), 0)
	if err != nil {
		return err
	}
	hashedSecret := string(hash[:])
	e.Secret = &hashedSecret

	log.Printf("Secret set for '%s'", e.GetID())
	return nil
}

// ValidateEntitySecretByID validates the identity of an entity by
// validating the authenticating entity with the secret.
func ValidateEntitySecretByID(ID string, secret string) error {
	e, err := GetEntityByID(ID)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(*e.Secret), []byte(secret))
	if err != nil {
		// This is strictly not in the style of go, but this
		// is the best place to put this log message so that
		// it works like all the others.
		log.Printf("Failed to authenticate '%s'", e.GetID())
	}
	log.Printf("Successfully authenticated '%s'", e.GetID())

	return err
}
