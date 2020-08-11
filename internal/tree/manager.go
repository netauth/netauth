package tree

import (
	"github.com/hashicorp/go-hclog"

	"github.com/netauth/netauth/internal/crypto"
	"github.com/netauth/netauth/internal/db"
)

// New returns an initialized tree.Manager on to which all other
// functions are bound.
func New(db db.DB, crypto crypto.EMCrypto, l hclog.Logger) (*Manager, error) {
	x := Manager{}
	x.log = l.Named("tree")
	x.bootstrapDone = false
	x.db = db
	x.crypto = crypto
	x.refContext = RefContext{
		DB:     db,
		Crypto: crypto,
	}

	// Initialize all entity hooks and bind to names.
	x.entityHooks = make(map[string]EntityHook)
	x.InitializeEntityHooks()

	// Construct entity chains out of the bound plugins.
	x.entityProcesses = make(map[string][]EntityHook)
	x.InitializeEntityChains(defaultEntityChains)

	// Check that required chains are loaded, bailing out if they
	// aren't.
	if err := x.CheckRequiredEntityChains(); err != nil {
		return nil, err
	}

	// Initialize all group hooks and bind to names.
	x.groupHooks = make(map[string]GroupHook)
	x.InitializeGroupHooks()

	// Construct group chains out of the bound plugins.
	x.groupProcesses = make(map[string][]GroupHook)
	x.InitializeGroupChains(defaultGroupChains)

	// Check that required chains are loaded, bailing out if they aren't.
	if err := x.CheckRequiredGroupChains(); err != nil {
		return nil, err
	}

	x.log.Debug("Initialized new Entity Manager")

	return &x, nil
}
