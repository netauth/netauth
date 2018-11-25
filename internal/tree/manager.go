package tree

import (
	"flag"
	"log"

	"github.com/NetAuth/NetAuth/internal/crypto"
	"github.com/NetAuth/NetAuth/internal/db"
)

var (
	debugChains = flag.Bool("verbose_chains", false, "Print verbose chain startup information")
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
	x.entityProcessorHooks = make(map[string]EntityHook)
	x.InitializeEntityHooks()

	// Construct entity chains out of the bound plugins.
	x.entityProcesses = make(map[string][]EntityHook)
	x.InitializeEntityChains(defaultEntityChains)

	// Check that required chains are loaded, bailing out if they
	// aren't.
	x.CheckRequiredEntityChains()

	// Initialize all group hooks and bind to names.
	x.groupProcessorHooks = make(map[string]GroupHook)
	x.InitializeGroupHooks()

	// Construct group chains out of the bound plugins.
	x.groupProcesses = make(map[string][]GroupHook)
	x.InitializeGroupChains(defaultGroupChains)

	// Check that required chains are loaded, bailing out if they aren't.
	x.CheckRequiredGroupChains()

	log.Println("Initialized new Entity Manager")

	return &x
}
