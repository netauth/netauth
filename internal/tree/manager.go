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
	x.processors = make(map[string]EntityProcessor)

	log.Println("Initialized new Entity Manager")

	return &x
}
