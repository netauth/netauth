package redis

import (
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/netauth/netauth/internal/db"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	testURL = "redis://127.0.0.1:1234"
)

func TestCB(t *testing.T) {
	cb()
}

func TestSetEventFunc(t *testing.T) {
	viper.Set("redis.url", testURL)
	kv, err := New(hclog.NewNullLogger())
	assert.Nil(t, err)

	f := func(db.Event) {}

	if kv.(*RedisStore).eF != nil {
		t.Log("EventFunc somehow already set!")
	}

	kv.SetEventFunc(f)

	if kv.(*RedisStore).eF == nil {
		t.Error("EventFunc not set correctly!")
	}
}

func TestClose(t *testing.T) {
	viper.Set("redis.url", testURL)
	kv, err := New(hclog.NewNullLogger())
	assert.Nil(t, err)
	kv.SetEventFunc(func(db.Event) {})

	assert.Nil(t, kv.Close())
}

func TestCapabilities(t *testing.T) {
	viper.Set("redis.url", testURL)
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
	viper.Set("redis.url", testURL)

	kv, err := New(hclog.NewNullLogger())
	assert.Nil(t, err)

	kv.SetEventFunc(ef.FireEvent)

	ef.On("FireEvent", db.Event{PK: "entity1", Type: db.EventEntityUpdate})
	ef.On("FireEvent", db.Event{PK: "entity1", Type: db.EventEntityDestroy})
	ef.On("FireEvent", db.Event{PK: "group1", Type: db.EventGroupUpdate})
	ef.On("FireEvent", db.Event{PK: "group1", Type: db.EventGroupDestroy})

	kv.(*RedisStore).fireEventForKey("/entities/entity1", eventUpdate)
	ef.AssertCalled(t, "FireEvent", db.Event{PK: "entity1", Type: db.EventEntityUpdate})

	kv.(*RedisStore).fireEventForKey("/entities/entity1", eventDelete)
	ef.AssertCalled(t, "FireEvent", db.Event{PK: "entity1", Type: db.EventEntityDestroy})

	kv.(*RedisStore).fireEventForKey("/groups/group1", eventUpdate)
	ef.AssertCalled(t, "FireEvent", db.Event{PK: "group1", Type: db.EventGroupUpdate})

	kv.(*RedisStore).fireEventForKey("/groups/group1", eventDelete)
	ef.AssertCalled(t, "FireEvent", db.Event{PK: "group1", Type: db.EventGroupDestroy})

	kv.(*RedisStore).fireEventForKey("/not/an/event/key", eventUpdate)

	// This assertion validates that the call for a key outside
	// the fixed keyspace went to the default case and didn't
	// trigger an event.
	ef.AssertNumberOfCalls(t, "FireEvent", 4)
}
