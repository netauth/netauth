package entity_manager

import (
	"log"

	pb "github.com/NetAuth/NetAuth/proto"
)

var (
	// e is a package scoped map of entities by string ID.
	eByID        = make(map[string]*pb.Entity)
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
	}

	eByID[ID] = newEntity
	eByUIDNumber[uidNumber] = newEntity

	// Successfully created, we now return no errors
	log.Printf("Created entity '%s'", ID)
	return nil
}

func GetEntityByID(ID string) (*pb.Entity, error) {
	e, ok := eByID[ID]
	if !ok {
		return nil, E_NO_ENTITY
	}
	return e, nil
}

func GetEntityByUIDNumber(uidNumber int32) (*pb.Entity, error) {
	e, ok := eByUIDNumber[uidNumber]
	if !ok {
		return nil, E_NO_ENTITY
	}
	return e, nil
}

func DeleteEntityByID(ID string) error {
	e, err := GetEntityByID(ID)
	if err != nil {
		return E_NO_ENTITY
	}

	// Now we need to delete from both maps
	delete(eByID, e.GetID())
	delete(eByUIDNumber, e.GetUidNumber())
	log.Printf("Deleted entity '%s'", e.GetID())

	// There's a small chance that the deletes above didn't go
	// through but if that happened there's not much that we can
	// do about it since delete() is a language construct and
	// there's no way to drill down deeper.  To be paranoid though
	// we'll clear the secret on the entity which should prevent
	// it from being usable.
	e.Secret = nil

	return nil
}

func DeleteEntityByUIDNumber(n int32) error {
	e, err := GetEntityByUIDNumber(n)
	if err != nil {
		return E_NO_ENTITY
	}

	// Now we need to delete from both maps
	delete(eByID, e.GetID())
	delete(eByUIDNumber, e.GetUidNumber())
	log.Printf("Deleted entity '%s'", e.GetID())

	// There's a small chance that the deletes above didn't go
	// through but if that happened there's not much that we can
	// do about it since delete() is a language construct and
	// there's no way to drill down deeper.  To be paranoid though
	// we'll clear the secret on the entity which should prevent
	// it from being usable.
	e.Secret = nil

	return nil
}
