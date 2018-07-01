package tree

import (
	"github.com/NetAuth/NetAuth/internal/crypto"
	"github.com/NetAuth/NetAuth/internal/db"
)

type Manager struct {
	// Making a bootstrap entity is a rare thing and short
	// circuits most of the permissions logic.  As such we only
	// allow it to be done once per server start.
	bootstrap_done bool

	// The persistence layer contains the functions that actually
	// deal with the disk and make this a useable server.
	db db.DB

	// The Crypto layer allows us to plug in different crypto
	// engines
	crypto crypto.EMCrypto
}
