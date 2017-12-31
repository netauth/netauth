package entity_tree

import "errors"

var (
	E_DUPLICATE_ID = errors.New("An entity with that ID already exists!")
)
