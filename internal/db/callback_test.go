package db

import (
	"testing"
)

var (
	dummyCalled = false
)

func dummyCallback(Event) { dummyCalled = true }

func TestCallbacks(t *testing.T) {
	RegisterCallback("foo", dummyCallback)
	RegisterCallback("foo", dummyCallback)

	if len(callbacks) != 1 {
		t.Error("Duplicate callback registered")
	}

	e := Event{
		Type: EventEntityCreate,
		PK:   "null",
	}

	FireEvent(e)

	if !dummyCalled {
		t.Error("Callbacks run but dummy was not called")
	}
	DeregisterCallback("foo")
}
