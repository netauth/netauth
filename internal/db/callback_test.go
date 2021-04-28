package db

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	dummyCalled = false
)

func dummyCallback(Event) { dummyCalled = true }

func TestCallbacks(t *testing.T) {
	x := &DB{cbs: make(map[string]Callback)}

	x.RegisterCallback("foo", dummyCallback)
	x.RegisterCallback("foo", dummyCallback)

	if len(x.cbs) != 1 {
		t.Error("Duplicate callback registered")
	}

	e := Event{
		Type: EventEntityCreate,
		PK:   "null",
	}

	x.FireEvent(e)

	if !dummyCalled {
		t.Error("Callbacks run but dummy was not called")
	}
}

func TestEventUpdateAll(t *testing.T) {
	c1Called := false
	c1 := func(e Event) {
		if e.Type == EventEntityUpdate {
			c1Called = true
		}
	}
	c2Called := false
	c2 := func(e Event) {
		if e.Type == EventGroupUpdate {
			c2Called = true
		}
	}

	RegisterKV("mock", newMockKV)
	x, _ := New("mock")
	delete(x.cbs, "BleveSearch") // Remove the search so we don't need to mock Get()

	entList := []string{"/entities/foo"}
	x.kv.(*mockKV).On("Keys", "/entities/*").Return(entList, nil)
	grpList := []string{"/groups/foo"}
	x.kv.(*mockKV).On("Keys", "/groups/*").Return(grpList, nil)

	x.RegisterCallback("c1", c1)
	x.RegisterCallback("c2", c2)

	assert.Nil(t, x.EventUpdateAll())

	if !(c1Called && c2Called) {
		t.Errorf("Not all callbacks satisfied; c1: %v, c2: %v", c1Called, c2Called)
	}
}

func TestEventAllBadEntity(t *testing.T) {
	RegisterKV("mock", newMockKV)
	x, _ := New("mock")
	delete(x.cbs, "BleveSearch") // Remove the search so we don't need to mock Get()

	x.kv.(*mockKV).On("Keys", "/entities/*").Return([]string{}, errors.New("entity error"))
	assert.NotNil(t, x.EventUpdateAll())
}

func TestEventAllBadGroup(t *testing.T) {
	RegisterKV("mock", newMockKV)
	x, _ := New("mock")
	delete(x.cbs, "BleveSearch") // Remove the search so we don't need to mock Get()

	x.kv.(*mockKV).On("Keys", "/entities/*").Return([]string{}, nil)
	x.kv.(*mockKV).On("Keys", "/groups/*").Return([]string{}, errors.New("group error"))
	assert.NotNil(t, x.EventUpdateAll())
}

func TestEventIsEmpty(t *testing.T) {
	e := Event{}

	if !e.IsEmpty() {
		t.Error("Empty event is claimed not empty!")
	}

	e.PK = "something"

	if e.IsEmpty() {
		t.Error("Filled event is claimed empty!")
	}
}
