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

	// ErrEntityLocked is returned when certain actions are
	// attempted on a locked entity.  Locked entities cannot
	// authenticate or change secrets.  They are effectively dead
	// to the system.
	ErrEntityLocked = errors.New("this entity is locked")

	// ErrHookExists is returned when a hook attempts to register
	// for a name that is already registered in the system.
	ErrHookExists = errors.New("a hook with this name already exists")

	// ErrUnknownHook is returned when a loader tries to add a
	// hook that is unknown to the chain.
	ErrUnknownHook = errors.New("no hook with this name exists")

	// ErrUnknownHookChain is returned when a processor attempts
	// to grab hooks from an unknown chain.
	ErrUnknownHookChain = errors.New("no chain with that ID exists")

	// ErrEmptyHookChain is returned when a chain was successfully
	// aquired, but it was empty.  In theory this shouldn't ever
	// happen, but its possible.
	ErrEmptyHookChain = errors.New("the specified chain is empty")
)
