package rpc2

import (
	"testing"

	"github.com/NetAuth/NetAuth/internal/crypto/nocrypto"
	"github.com/NetAuth/NetAuth/internal/db/memdb"
	"github.com/NetAuth/NetAuth/internal/token/null"
	"github.com/NetAuth/NetAuth/internal/tree"
	_ "github.com/NetAuth/NetAuth/internal/tree/hooks"

	types "github.com/NetAuth/Protocol"
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

func initTree(t *testing.T, m Manager) {
	m.CreateEntity("admin", -1, "secret")
	m.CreateEntity("unprivileged", -1, "secret")
	m.CreateEntity("entity1", -1, "secret")

	m.CreateGroup("group1", "", "", -1)
	m.CreateGroup("group2", "", "group1", -1)

	m.AddEntityToGroup("entity1", "group1")

	m.SetEntityCapability2("admin", types.Capability_GLOBAL_ROOT.Enum())
}

func TestNew(t *testing.T) {
	// This is the most basic check to make sure that all the shim
	// interfaces generate correctly.  This test should fail first
	// and most obviously if that's the case.
	newServer(t)
}
