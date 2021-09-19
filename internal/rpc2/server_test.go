package rpc2

import (
	"path"
	"testing"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/protobuf/proto"

	"github.com/netauth/netauth/internal/crypto/nocrypto"
	"github.com/netauth/netauth/internal/db"
	"github.com/netauth/netauth/internal/db/memory"
	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/token/null"
	"github.com/netauth/netauth/internal/tree"
	_ "github.com/netauth/netauth/internal/tree/hooks"

	types "github.com/netauth/protocol"
)

type errorableKV struct {
	db.KVStore
}

func (e *errorableKV) Put(k string, v []byte) error {
	if path.Base(k) == "save-error" {
		return db.ErrInternalError
	}
	return e.KVStore.Put(k, v)
}

func (e *errorableKV) Get(k string) ([]byte, error) {
	if path.Base(k) == "load-error" {
		return nil, db.ErrInternalError
	}
	return e.KVStore.Get(k)
}

func newServer(t *testing.T) *Server {
	startup.DoCallbacks()

	db.RegisterKV("errorable", func(l hclog.Logger) (db.KVStore, error) {
		mkv, _ := memory.NewKV(l)
		return &errorableKV{mkv}, nil
	})

	db, err := db.New("errorable")
	if err != nil {
		t.Fatal(err)
	}

	crypto, err := nocrypto.New(hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}

	m, err := tree.New(db, crypto, hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}

	n := null.New(hclog.NewNullLogger())

	return New(Refs{TokenService: n, Tree: m}, hclog.NewNullLogger())
}

func newServerWithRefs(t *testing.T) (*Server, tree.DB, Manager) {
	startup.DoCallbacks()

	db.RegisterKV("errorable", func(l hclog.Logger) (db.KVStore, error) {
		mkv, _ := memory.NewKV(l)
		return &errorableKV{mkv}, nil
	})

	db, err := db.New("errorable")
	if err != nil {
		t.Fatal(err)
	}

	crypto, err := nocrypto.New(hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}

	m, err := tree.New(db, crypto, hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}

	n := null.New(hclog.NewNullLogger())

	return New(Refs{TokenService: n, Tree: m}, hclog.NewNullLogger()), db, m
}

func initTree(t *testing.T, m Manager) {
	m.CreateEntity("admin", -1, "secret")
	m.CreateEntity("unprivileged", -1, "secret")
	m.CreateEntity("entity1", -1, "secret")

	m.CreateGroup("group1", "", "", -1)
	m.CreateGroup("group2", "", "group1", -1)

	m.AddEntityToGroup("entity1", "group1")

	m.SetEntityCapability2("admin", types.Capability_GLOBAL_ROOT.Enum())

	m.GroupKVAdd("group1", []*types.KVData{{Key: proto.String("key1"), Values: []*types.KVValue{{Value: proto.String("value1")}}}})
	m.EntityKVAdd("entity1", []*types.KVData{{Key: proto.String("key1"), Values: []*types.KVValue{{Value: proto.String("value1")}}}})
}

func TestNew(t *testing.T) {
	// This is the most basic check to make sure that all the shim
	// interfaces generate correctly.  This test should fail first
	// and most obviously if that's the case.
	newServer(t)
}
