package entity_manager

import "errors"

var (
	E_DUPLICATE_ID        = errors.New("An entity with that ID already exists!")
	E_DUPLICATE_UIDNUMBER = errors.New("An entity with that uidNumber already exists!")
	E_NO_ENTITY = errors.New("No entity matched the given constraints")
)
