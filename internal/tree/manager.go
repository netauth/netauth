package tree

import (
	"github.com/hashicorp/go-hclog"

	"github.com/netauth/netauth/internal/crypto"
	"github.com/netauth/netauth/internal/mresolver"
)

var (
	// This logger is very special, and is meant for exclusive use
	// by early init tasks that happen at the package scope and
	// are not bound to a tree instance.
	initlb hclog.Logger
)

// New returns an initialized tree.Manager on to which all other
// functions are bound.
func New(db DB, crypto crypto.EMCrypto, l hclog.Logger) (*Manager, error) {
	x := Manager{}
	x.log = l.Named("tree")
	x.db = db
	x.crypto = crypto
	x.resolver = mresolver.New()
	x.resolver.SetParentLogger(x.log)
	x.refContext = RefContext{
		DB:     db,
		Crypto: crypto,
	}

	x.db.RegisterCallback("entity-resolver", x.entityResolverCallback)
	x.db.RegisterCallback("group-resolver", x.groupResolverCallback)

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

// SetParentLogger sets the parent logger for this instance.
func SetParentLogger(l hclog.Logger) {
	initlb = l.Named("tree.init")
}

// log is a convenience function that will return a null logger if a
// parent logger has not been specified, mostly useful for tests.
func log() hclog.Logger {
	if initlb == nil {
		initlb = hclog.NewNullLogger()
	}
	return initlb
}
