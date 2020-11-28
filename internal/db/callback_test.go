package db

import (
	"testing"
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
