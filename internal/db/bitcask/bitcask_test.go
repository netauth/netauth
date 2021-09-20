package bitcask

import (
	"context"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/netauth/netauth/internal/db"
)

func TestCB(t *testing.T) {
	cb()
}

func TestNewBadLock(t *testing.T) {
	viper.Set("core.home", t.TempDir())

	// This one should work
	_, err := New(hclog.NewNullLogger())
	assert.Nil(t, err)

	// This one shouldn't
	_, err = New(hclog.NewNullLogger())
	assert.NotNil(t, err)
}

func TestSetEventFunc(t *testing.T) {
	viper.Set("core.home", t.TempDir())
	kv, err := New(hclog.NewNullLogger())
	assert.Nil(t, err)

	f := func(db.Event) {}

	if kv.(*BCStore).eF != nil {
		t.Log("EventFunc somehow already set!")
	}

	kv.SetEventFunc(f)

	if kv.(*BCStore).eF == nil {
		t.Error("EventFunc not set correctly!")
	}
}

func TestPut(t *testing.T) {
	ctx := context.Background()
	viper.Set("core.home", t.TempDir())
	kv, err := New(hclog.NewNullLogger())
	assert.Nil(t, err)
	kv.SetEventFunc(func(db.Event) {})

	assert.Nil(t, kv.Put(ctx, "/entities/entity1", []byte("some data")))
	assert.Nil(t, kv.Put(ctx, "/groups/group1", []byte("some more data")))

	b, err := kv.Get(ctx, "/entities/entity1")
	assert.Nil(t, err)
	assert.Equal(t, []byte("some data"), b, "Data stored but incorrect")

	assert.NotNil(t, kv.Put(ctx, "", []byte("data-with-no-key")))
}

func TestGet(t *testing.T) {
	ctx := context.Background()
	viper.Set("core.home", t.TempDir())
	kv, err := New(hclog.NewNullLogger())
	assert.Nil(t, err)
	kv.SetEventFunc(func(db.Event) {})

	assert.Nil(t, kv.Put(ctx, "/entities/entity1", []byte("lots of data")))

	v, err := kv.Get(ctx, "/entities/entity1")
	assert.Nil(t, err)
	assert.Equal(t, v, []byte("lots of data"), "Data read but incorrect")

	v, err = kv.Get(ctx, "/does/not/exist")
	assert.Nil(t, v, "Data returned for key that does not exist")
	assert.Equal(t, err, db.ErrNoValue, "KV made up some data")
}

func TestDel(t *testing.T) {
	ctx := context.Background()
	viper.Set("core.home", t.TempDir())
	kv, err := New(hclog.NewNullLogger())
	assert.Nil(t, err)
	kv.SetEventFunc(func(db.Event) {})

	assert.Nil(t, kv.Put(ctx, "/entities/entity1", []byte("lots of data")))
	assert.Nil(t, kv.Put(ctx, "/groups/group1", []byte("lots of data")))

	_, err = kv.Get(ctx, "/entities/entity1")
	assert.Nil(t, err)
	_, err = kv.Get(ctx, "/groups/group1")
	assert.Nil(t, err)

	assert.Nil(t, kv.Del(ctx, "/entities/entity1"))
	assert.Nil(t, kv.Del(ctx, "/groups/group1"))

	_, err = kv.Get(ctx, "/entities/entity1")
	assert.Equal(t, err, db.ErrNoValue)
	_, err = kv.Get(ctx, "/groups/group1")
	assert.Equal(t, err, db.ErrNoValue)
}

func TestKeys(t *testing.T) {
	ctx := context.Background()
	viper.Set("core.home", t.TempDir())
	kv, err := New(hclog.NewNullLogger())
	assert.Nil(t, err)
	kv.SetEventFunc(func(db.Event) {})

	kv.Put(ctx, "/entities/entity1", []byte("lots of data"))
	kv.Put(ctx, "/entities/entity2", []byte("lots of data"))
	kv.Put(ctx, "/entities/purple", []byte("lots of data"))

	res, err := kv.Keys(ctx, "/entities/entity*")
	assert.Nil(t, err)
	assert.Equal(t, res, []string{"/entities/entity1", "/entities/entity2"})
}

func TestClose(t *testing.T) {
	viper.Set("core.home", t.TempDir())
	kv, err := New(hclog.NewNullLogger())
	assert.Nil(t, err)
	kv.SetEventFunc(func(db.Event) {})

	assert.Nil(t, kv.Close())
}

func TestCapabilities(t *testing.T) {
	viper.Set("core.home", t.TempDir())
	kv, err := New(hclog.NewNullLogger())
	assert.Nil(t, err)
	kv.SetEventFunc(func(db.Event) {})

	assert.Equal(t, []db.KVCapability{db.KVMutable}, kv.Capabilities())
}

type eventHandler struct{ mock.Mock }

func (eh *eventHandler) FireEvent(e db.Event) {
	eh.Called(e)
}

func TestFireEventForKey(t *testing.T) {
	ef := eventHandler{}

	kv, err := New(hclog.NewNullLogger())
	assert.Nil(t, err)

	kv.SetEventFunc(ef.FireEvent)

	ef.On("FireEvent", db.Event{PK: "entity1", Type: db.EventEntityUpdate})
	ef.On("FireEvent", db.Event{PK: "entity1", Type: db.EventEntityDestroy})
	ef.On("FireEvent", db.Event{PK: "group1", Type: db.EventGroupUpdate})
	ef.On("FireEvent", db.Event{PK: "group1", Type: db.EventGroupDestroy})

	kv.(*BCStore).fireEventForKey("/entities/entity1", eventUpdate)
	ef.AssertCalled(t, "FireEvent", db.Event{PK: "entity1", Type: db.EventEntityUpdate})

	kv.(*BCStore).fireEventForKey("/entities/entity1", eventDelete)
	ef.AssertCalled(t, "FireEvent", db.Event{PK: "entity1", Type: db.EventEntityDestroy})

	kv.(*BCStore).fireEventForKey("/groups/group1", eventUpdate)
	ef.AssertCalled(t, "FireEvent", db.Event{PK: "group1", Type: db.EventGroupUpdate})

	kv.(*BCStore).fireEventForKey("/groups/group1", eventDelete)
	ef.AssertCalled(t, "FireEvent", db.Event{PK: "group1", Type: db.EventGroupDestroy})

	kv.(*BCStore).fireEventForKey("/not/an/event/key", eventUpdate)

	// This assertion validates that the call for a key outside
	// the fixed keyspace went to the default case and didn't
	// trigger an event.
	ef.AssertNumberOfCalls(t, "FireEvent", 4)
}
