package tree

import (
	"log"

	"github.com/NetAuth/NetAuth/internal/crypto"
	"github.com/NetAuth/NetAuth/internal/db"
)

// New returns an initialized tree.Manager on to which all other
// functions are bound.
func New(db db.DB, crypto crypto.EMCrypto) *Manager {
	x := Manager{}
	x.bootstrapDone = false
	x.db = db
	x.crypto = crypto
	x.refContext = RefContext{
		DB:     db,
		Crypto: crypto,
	}

	// Initialize all entity hooks and bind to names.
	x.entityProcessorHooks = make(map[string]EntityProcessorHook)
	x.InitializeEntityHooks()

	// Construct entity chains out of the bound plugins.
	x.entityProcesses = make(map[string][]EntityProcessorHook)
	x.InitializeEntityChains(defaultEntityChains)

	// Initialize all group hooks and bind to names.
	x.groupProcessorHooks = make(map[string]GroupProcessorHook)
	x.InitializeGroupHooks()

	// Construct group chains out of the bound plugins.
	x.groupProcesses = make(map[string][]GroupProcessorHook)
	x.InitializeGroupChains(defaultGroupChains)

	log.Println("Initialized new Entity Manager")

	return &x
}
