package tree

import (
	"github.com/NetAuth/NetAuth/internal/crypto"
	"github.com/NetAuth/NetAuth/internal/db"
)

// The Manager binds all methods for managing a tree of entities with
// the associated groups, capabilities, and other assorted functions.
// This is the type that is served up by the RPC layer.
type Manager struct {
	// Making a bootstrap entity is a rare thing and short
	// circuits most of the permissions logic.  As such we only
	// allow it to be done once per server start.
	bootstrapDone bool

	// The persistence layer contains the functions that actually
	// deal with the disk and make this a useable server.
	db db.DB

	// The Crypto layer allows us to plug in different crypto
	// engines
	crypto crypto.EMCrypto

	// Maintain maps of hooks that have been initialized.
	entityProcessorHooks map[string]EntityProcessorHook
	groupProcessorHooks  map[string]GroupProcessorHook

	// Maintain chains of hooks that can be used by processors.
	entityProcesses map[string][]EntityProcessorHook
	groupProcesses  map[string][]GroupProcessorHook
}

// A RefContext is a container of references that are needed to
// bootstrap the tree manager and associated plugins.
type RefContext struct {
	DB     db.DB
	crypto crypto.EMCrypto
}
