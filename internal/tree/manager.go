package tree

import (
	"log"

	"github.com/NetAuth/NetAuth/internal/crypto"
	"github.com/NetAuth/NetAuth/internal/db"

	"github.com/NetAuth/NetAuth/internal/tree/hooks"
)

// New returns an initialized tree.Manager on to which all other
// functions are bound.
func New(db db.DB, crypto crypto.EMCrypto) *Manager {
	x := Manager{}
	x.bootstrapDone = false
	x.db = db
	x.crypto = crypto
	x.processors = make(map[string]EntityProcessor)

	x.entityProcesses = make(map[string][]EntityProcessorHook)

	log.Println("Initialized new Entity Manager")

	x.EntityHookMustRegister("CREATE", &hooks.FailOnExistingEntity{db})
	x.EntityHookMustRegister("CREATE", &hooks.SetEntityID{})
	x.EntityHookMustRegister("CREATE", &hooks.SetEntityNumber{db})
	x.EntityHookMustRegister("CREATE", &hooks.SetEntitySecret{crypto})
	x.EntityHookMustRegister("CREATE", &hooks.SaveEntity{db})

	x.EntityHookMustRegister("BOOTSTRAP-SERVER", &hooks.CreateEntityIfMissing{db})
	x.EntityHookMustRegister("BOOTSTRAP-SERVER", &hooks.EnsureEntityMeta{})
	x.EntityHookMustRegister("BOOTSTRAP-SERVER", &hooks.UnlockEntity{})
	x.EntityHookMustRegister("BOOTSTRAP-SERVER", &hooks.SetEntityCapability{})
	x.EntityHookMustRegister("BOOTSTRAP-SERVER", &hooks.SaveEntity{db})

	x.EntityHookMustRegister("DESTROY", &hooks.DestroyEntity{db})

	x.EntityHookMustRegister("SET-CAPABILITY", &hooks.LoadEntity{db})
	x.EntityHookMustRegister("SET-CAPABILITY", &hooks.EnsureEntityMeta{})
	x.EntityHookMustRegister("SET-CAPABILITY", &hooks.SetEntityCapability{})
	x.EntityHookMustRegister("SET-CAPABILITY", &hooks.SaveEntity{db})

	x.EntityHookMustRegister("DROP-CAPABILITY", &hooks.LoadEntity{db})
	x.EntityHookMustRegister("DROP-CAPABILITY", &hooks.EnsureEntityMeta{})
	x.EntityHookMustRegister("DROP-CAPABILITY", &hooks.RemoveEntityCapability{})
	x.EntityHookMustRegister("DROP-CAPABILITY", &hooks.SaveEntity{db})

	x.EntityHookMustRegister("SET-SECRET", &hooks.LoadEntity{db})
	x.EntityHookMustRegister("SET-SECRET", &hooks.SetEntitySecret{crypto})
	x.EntityHookMustRegister("SET-SECRET", &hooks.SaveEntity{db})

	x.EntityHookMustRegister("VALIDATE-IDENTITY", &hooks.LoadEntity{db})
	x.EntityHookMustRegister("VALIDATE-IDENTITY", &hooks.ValidateEntityUnlocked{})
	x.EntityHookMustRegister("VALIDATE-IDENTITY", &hooks.ValidateEntitySecret{crypto})
	x.EntityHookMustRegister("VALIDATE-IDENTITY", &hooks.SaveEntity{db})

	x.EntityHookMustRegister("FETCH", &hooks.LoadEntity{db})

	x.EntityHookMustRegister("MERGE-METADATA", &hooks.LoadEntity{db})
	x.EntityHookMustRegister("MERGE-METADATA", &hooks.EnsureEntityMeta{})
	x.EntityHookMustRegister("MERGE-METADATA", &hooks.MergeEntityMeta{})
	x.EntityHookMustRegister("MERGE-METADATA", &hooks.SaveEntity{db})

	x.EntityHookMustRegister("LOCK", &hooks.LoadEntity{db})
	x.EntityHookMustRegister("LOCK", &hooks.EnsureEntityMeta{})
	x.EntityHookMustRegister("LOCK", &hooks.LockEntity{})
	x.EntityHookMustRegister("LOCK", &hooks.SaveEntity{db})

	x.EntityHookMustRegister("UNLOCK", &hooks.LoadEntity{db})
	x.EntityHookMustRegister("UNLOCK", &hooks.EnsureEntityMeta{})
	x.EntityHookMustRegister("UNLOCK", &hooks.UnlockEntity{})
	x.EntityHookMustRegister("UNLOCK", &hooks.SaveEntity{db})

	x.EntityHookMustRegister("ADD-KEY", &hooks.LoadEntity{db})
	x.EntityHookMustRegister("ADD-KEY", &hooks.EnsureEntityMeta{})
	x.EntityHookMustRegister("ADD-KEY", &hooks.AddEntityKey{})
	x.EntityHookMustRegister("ADD-KEY", &hooks.SaveEntity{db})

	x.EntityHookMustRegister("DEL-KEY", &hooks.LoadEntity{db})
	x.EntityHookMustRegister("DEL-KEY", &hooks.EnsureEntityMeta{})
	x.EntityHookMustRegister("DEL-KEY", &hooks.DelEntityKey{})
	x.EntityHookMustRegister("DEL-KEY", &hooks.SaveEntity{db})
	
	return &x
}
