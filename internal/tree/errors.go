package tree

import "errors"

var (
	DuplicateEntityID  = errors.New("This ID is already allocated")
	DuplicateGroupName = errors.New("This name is already allocated")
	DuplicateNumber    = errors.New("This number is already allocated")
	UnknownCapability  = errors.New("The capability specified is unknown")
	ExistingExpansion  = errors.New("This expansion already exists!")
)
