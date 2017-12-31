package entity_tree

import (
	"log"

	pb "github.com/NetAuth/NetAuth/proto"
)

var (
	// e is a package scoped map of entities by string ID.
	e = make(map[string]*pb.Entity)
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
	for i := range e {
		if e[i].GetUidNumber() > largest {
			largest = e[i].GetUidNumber()
		}
	}

	return largest + 1
}

func NewEntity(ID string, uidNumber int32, secret string) error {
	// Does this entity exist already?
	_, exists := e[ID]
	if exists {
		log.Printf("Entity with ID '%s' already exists!", ID)
		return E_DUPLICATE_ID
	}

	// Were we given a specific uidNumber?
	if uidNumber == -1 {
		// -1 is a sentinel value that tells us to pick the
		// next available number and assign it.
		uidNumber = nextUIDNumber()
	}

	// Ok, they don't exist so we'll make them exist now
	e[ID] = &pb.Entity{
		ID:        &ID,
		UidNumber: &uidNumber,
		Secret:    &secret,
	}

	// Successfully created, we now return no errors
	log.Printf("Created entity '%s'", ID)
	return nil
}
