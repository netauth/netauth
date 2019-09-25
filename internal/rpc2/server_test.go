package rpc2

import (
	"testing"

	"github.com/NetAuth/NetAuth/internal/crypto/nocrypto"
	"github.com/NetAuth/NetAuth/internal/db/memdb"
	"github.com/NetAuth/NetAuth/internal/tree"
	_ "github.com/NetAuth/NetAuth/internal/tree/hooks"
	"github.com/NetAuth/NetAuth/internal/token/null"
)

func newServer(t *testing.T) *Server {
	db, err := memdb.New()
	if err != nil {
		t.Fatal(err)
	}

	crypto, err := nocrypto.New()
	if err != nil {
		t.Fatal(err)
	}

	m, err := tree.New(db, crypto)
	if err != nil {
		t.Fatal(err)
	}

	n := null.New()

	return New(Refs{TokenService: n, Tree: m})
}

func TestNew(t *testing.T) {
	// This is the most basic check to make sure that all the shim
	// interfaces generate correctly.  This test should fail first
	// and most obviously if that's the case.
	newServer(t)
}
