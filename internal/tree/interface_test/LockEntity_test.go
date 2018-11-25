package interface_test

import (
	"testing"
)

func TestLockEntity(t *testing.T) {
	m, ctx := newTreeManager(t)

	addEntity(t, ctx)

	if err := m.LockEntity("entity1"); err != nil {
		t.Fatal(err)
	}

	e, err := ctx.DB.LoadEntity("entity1")
	if err != nil {
		t.Fatal(err)
	}

	if !e.GetMeta().GetLocked() {
		t.Error("Entity not locked")
	}

	if err := m.UnlockEntity("entity1"); err != nil {
		t.Fatal(err)
	}

	e, err = ctx.DB.LoadEntity("entity1")
	if err != nil {
		t.Fatal(err)
	}

	if e.GetMeta().GetLocked() {
		t.Error("Entity not unlocked")
	}	
}
