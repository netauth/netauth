package tree

import "errors"

var (
	DuplicateEntityID    = errors.New("This ID is already allocated")
	DuplicateGroupName   = errors.New("This name is already allocated")
	DuplicateNumber      = errors.New("This number is already allocated")
	UnknownCapability    = errors.New("The capability specified is unknown")
	AuthorizationFailure = errors.New("Authorization failed - bad credentials")
	ExistingExpansion    = errors.New("This expansion already exists!")
)
