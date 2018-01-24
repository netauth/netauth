package MemDB

// This is an incredibly simple testing file since the MemDB is a shim
// implementation

import (
	"testing"

	"github.com/NetAuth/NetAuth/internal/server/entity_manager"
)

func TestDiscoverEntityIDs(t *testing.T) {
	x := New()

	l, err := x.DiscoverEntityIDs()
	if len(l) != 0 {
		t.Error("MemDB invented entities!")
	}
	if err != nil {
		t.Error("MemDB made up an error!?")
	}
}

func TestLoadEntity(t *testing.T) {
	x := New()

	if _, err := x.LoadEntity(""); err != entity_manager.E_NO_ENTITY {
		t.Errorf("MemDB unexpected error: %s", err)
	}
}

func TestSaveEntity(t *testing.T) {
	x := New()

	if x.SaveEntity(nil) != nil {
		t.Error("MemDB made up an error!?")
	}
}

func TestDeleteEntity(t *testing.T) {
	x := New()

	if x.DeleteEntity("") != nil {
		t.Error("MemDB made up an error!?")
	}
}
