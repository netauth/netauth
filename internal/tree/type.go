package tree

import (
	"github.com/hashicorp/go-hclog"

	"github.com/netauth/netauth/internal/crypto"
	"github.com/netauth/netauth/internal/db"
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

	// A refContext maintains pointers to all referenced
	// subsystems required by a tree manager.
	refContext RefContext

	// Maintain maps of hooks that have been initialized.
	entityHooks map[string]EntityHook
	groupHooks  map[string]GroupHook

	// Maintain chains of hooks that can be used by processors.
	entityProcesses map[string][]EntityHook
	groupProcesses  map[string][]GroupHook

	log hclog.Logger
}

// A RefContext is a container of references that are needed to
// bootstrap the tree manager and associated plugins.
type RefContext struct {
	DB     db.DB
	Crypto crypto.EMCrypto
}

// The ChainConfig type maps from chain name to a list of hooks that
// should be in this chain.  The same type is used for entities and
// groups, but as these each have separate chains, different configs
// must be created and loaded for each.
type ChainConfig map[string][]string
