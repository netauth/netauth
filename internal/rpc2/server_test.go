package rpc2

import (
	"context"
	"path"
	"testing"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/protobuf/proto"

	"github.com/netauth/netauth/internal/crypto/nocrypto"
	"github.com/netauth/netauth/internal/db"
	"github.com/netauth/netauth/internal/db/memory"
	"github.com/netauth/netauth/internal/startup"
	"github.com/netauth/netauth/internal/tree"
	_ "github.com/netauth/netauth/internal/tree/hooks"
	"github.com/netauth/netauth/pkg/token/null"

	types "github.com/netauth/protocol"
)

type errorableKV struct {
	db.KVStore
}

func (e *errorableKV) Put(ctx context.Context, k string, v []byte) error {
	if path.Base(k) == "save-error" {
		return db.ErrInternalError
	}
	return e.KVStore.Put(ctx, k, v)
}

func (e *errorableKV) Get(ctx context.Context, k string) ([]byte, error) {
	if path.Base(k) == "load-error" {
		return nil, db.ErrInternalError
	}
	return e.KVStore.Get(ctx, k)
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

	m, err := tree.New(tree.WithStorage(db), tree.WithCrypto(crypto))
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

	m, err := tree.New(tree.WithStorage(db), tree.WithCrypto(crypto))
	if err != nil {
		t.Fatal(err)
	}

	n := null.New(hclog.NewNullLogger())

	return New(Refs{TokenService: n, Tree: m}, hclog.NewNullLogger()), db, m
}

func initTree(t *testing.T, m Manager) {
	ctx := context.Background()
	m.CreateEntity(ctx, "admin", -1, "secret")
	m.CreateEntity(ctx, "unprivileged", -1, "secret")
	m.CreateEntity(ctx, "entity1", -1, "secret")

	m.CreateGroup(ctx, "group1", "", "", -1)
	m.CreateGroup(ctx, "group2", "", "group1", -1)

	m.AddEntityToGroup(ctx, "entity1", "group1")

	m.SetEntityCapability2(ctx, "admin", types.Capability_GLOBAL_ROOT.Enum())

	m.GroupKVAdd(ctx, "group1", []*types.KVData{{Key: proto.String("key1"), Values: []*types.KVValue{{Value: proto.String("value1")}}}})
	m.EntityKVAdd(ctx, "entity1", []*types.KVData{{Key: proto.String("key1"), Values: []*types.KVValue{{Value: proto.String("value1")}}}})
}

func TestNew(t *testing.T) {
	// This is the most basic check to make sure that all the shim
	// interfaces generate correctly.  This test should fail first
	// and most obviously if that's the case.
	newServer(t)
}
