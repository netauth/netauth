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

	x.entityProcessorHooks = make(map[string]EntityProcessorHook)
	x.InitializeEntityHooks()

	x.entityProcesses = make(map[string][]EntityProcessorHook)
	x.InitializeEntityChains(defaultEntityChains)

	log.Println("Initialized new Entity Manager")

	// // Now the groups
	// x.groupProcesses = make(map[string][]GroupProcessorHook)

	// x.GroupHookMustRegister("CREATE-GROUP", &hooks.FailOnExistingGroup{db})
	// x.GroupHookMustRegister("CREATE-GROUP", &hooks.SetGroupName{})
	// x.GroupHookMustRegister("CREATE-GROUP", &hooks.SetManagingGroup{db})
	// x.GroupHookMustRegister("CREATE-GROUP", &hooks.SetGroupDisplayName{})
	// x.GroupHookMustRegister("CREATE-GROUP", &hooks.SetGroupNumber{db})
	// x.GroupHookMustRegister("CREATE-GROUP", &hooks.SaveGroup{db})

	return &x
}
