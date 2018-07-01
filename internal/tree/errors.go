package tree

import "errors"

var (
	// ErrDuplicateEntityID is returned when the entity ID
	// requested is already in use.
	ErrDuplicateEntityID = errors.New("this ID is already allocated")

	// ErrDuplicateGroupName is returned when the group name
	// requested is already in use.
	ErrDuplicateGroupName = errors.New("this name is already allocated")

	// ErrDuplicateNumber is returned if the number requested is
	// already in use.
	ErrDuplicateNumber = errors.New("this number is already allocated")

	// ErrUnknownCapability is returned when an action is
	// requested that involves a capability not known to the
	// system.
	ErrUnknownCapability = errors.New("the capability specified is unknown")

	// ErrExistingExpansion is returned when an action would
	// create an expansion that already exists.
	ErrExistingExpansion = errors.New("this expansion already exists")
)
