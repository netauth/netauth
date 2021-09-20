package memory

import (
	"context"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"

	"github.com/netauth/netauth/internal/db"
)

func TestCB(t *testing.T) {
	cb()
}

func TestSetEventFunc(t *testing.T) {
	kv, _ := NewKV(hclog.NewNullLogger())

	f := func(db.Event) {}

	if kv.(*KV).eF != nil {
		t.Log("EventFunc somehow already set!")
	}

	kv.SetEventFunc(f)

	if kv.(*KV).eF == nil {
		t.Error("EventFunc not set correctly!")
	}
}

func TestPut(t *testing.T) {
	ctx := context.Background()
	kv, _ := NewKV(hclog.NewNullLogger())
	kv.SetEventFunc(func(db.Event) {})

	assert.Nil(t, kv.Put(ctx, "/entities/entity1", []byte("some data")))
	assert.Nil(t, kv.Put(ctx, "/groups/group1", []byte("some more data")))

	assert.Equal(t, []byte("some data"), kv.(*KV).m["/entities/entity1"], "Data stored but incorrect")
}

func TestGet(t *testing.T) {
	ctx := context.Background()
	kv, _ := NewKV(hclog.NewNullLogger())
	kv.SetEventFunc(func(db.Event) {})

	kv.(*KV).m["/entities/entity1"] = []byte("lots of data")

	v, err := kv.Get(ctx, "/entities/entity1")
	assert.Nil(t, err)
	assert.Equal(t, v, []byte("lots of data"), "Data read but incorrect")

	v, err = kv.Get(ctx, "/does/not/exist")
	assert.Nil(t, v, "Data returned for key that does not exist")
	assert.Equal(t, err, db.ErrNoValue, "KV made up some data")
}

func TestDel(t *testing.T) {
	ctx := context.Background()
	kv, _ := NewKV(hclog.NewNullLogger())
	kv.SetEventFunc(func(db.Event) {})

	kv.(*KV).m["/entities/entity1"] = []byte("lots of data")
	kv.(*KV).m["/groups/group1"] = []byte("lots of data")

	assert.Nil(t, kv.Del(ctx, "/entities/entity1"))
	assert.Nil(t, kv.Del(ctx, "/groups/group1"))

	_, exists := kv.(*KV).m["/entities/entity1"]
	assert.False(t, exists, "Delete failed to remove entity")
	_, exists = kv.(*KV).m["/groups/group1"]
	assert.False(t, exists, "Delete failed to remove group")
}

func TestKeys(t *testing.T) {
	ctx := context.Background()
	kv, _ := NewKV(hclog.NewNullLogger())
	kv.SetEventFunc(func(db.Event) {})

	kv.(*KV).m["/entities/entity1"] = []byte("lots of data")
	kv.(*KV).m["/entities/entity2"] = []byte("lots of data")
	kv.(*KV).m["/entities/purple"] = []byte("lots of data")

	res, err := kv.Keys(ctx, "/entities/entity*")
	assert.Nil(t, err)
	assert.Equal(t, res, []string{"/entities/entity1", "/entities/entity2"})
}

func TestClose(t *testing.T) {
	kv, _ := NewKV(hclog.NewNullLogger())
	kv.SetEventFunc(func(db.Event) {})

	assert.Nil(t, kv.Close())
}

func TestCapabilities(t *testing.T) {
	kv, _ := NewKV(hclog.NewNullLogger())
	kv.SetEventFunc(func(db.Event) {})

	assert.Equal(t, []db.KVCapability{db.KVMutable}, kv.Capabilities())
}
