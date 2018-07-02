package tree

import (
	"testing"

	"github.com/NetAuth/NetAuth/internal/crypto/nocrypto"
	"github.com/NetAuth/NetAuth/internal/db/memdb"
)

func getNewEntityManager(t *testing.T) *Manager {
	db, err := memdb.New()
	if err != nil {
		t.Fatal(err)
	}

	crypto := nocrypto.New()
	return New(db, crypto)
}
