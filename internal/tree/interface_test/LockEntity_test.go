package interface_test

import (
	"context"
	"testing"
)

func TestLockEntity(t *testing.T) {
	ctxt := context.Background()
	m, mdb := newTreeManager(t)

	addEntity(t, mdb)

	if err := m.LockEntity(ctxt, "entity1"); err != nil {
		t.Fatal(err)
	}

	e, err := mdb.LoadEntity(ctxt, "entity1")
	if err != nil {
		t.Fatal(err)
	}

	if !e.GetMeta().GetLocked() {
		t.Error("Entity not locked")
	}

	if err := m.UnlockEntity(ctxt, "entity1"); err != nil {
		t.Fatal(err)
	}

	e, err = mdb.LoadEntity(ctxt, "entity1")
	if err != nil {
		t.Fatal(err)
	}

	if e.GetMeta().GetLocked() {
		t.Error("Entity not unlocked")
	}
}
