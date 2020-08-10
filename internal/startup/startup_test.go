package startup

import (
	"testing"
)

func TestRegisterCallback(t *testing.T) {
	callbacks = nil
	RegisterCallback(func() {})
	if len(callbacks) != 1 {
		t.Error("Callback not registered")
	}
}

func TestDoCallbacks(t *testing.T) {
	callbacks = nil
	called := false

	testCB := func() {
		called = true
	}

	RegisterCallback(testCB)
	DoCallbacks()

	if !called {
		t.Error("Callback was not called")
	}
}
