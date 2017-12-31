package entity_tree

import (
	"log"

	pb "github.com/NetAuth/NetAuth/proto"
)

var (
	// e is a package scoped map of entities by string ID.
	e = make(map[string]*pb.Entity)
)

func NewEntity(ID string, uidNumber int32, secret string) error {
	// Does this entity exist already?
	_, exists := e[ID]
	if exists {
		log.Printf("Entity with ID '%s' already exists!", ID)
		return E_DUPLICATE_ID
	}

	// Ok, they don't exist so we'll make them exist now

	e[ID] = &pb.Entity{
		ID: &ID,
		UidNumber: &uidNumber,
		Secret: &secret,
	}

	// Successfully created, we now return no errors
	log.Printf("Created entity '%s'", ID)
	return nil
}
