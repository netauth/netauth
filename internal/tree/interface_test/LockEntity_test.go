package interface_test

import (
	"context"
	"testing"
)

func TestLockEntity(t *testing.T) {
	ctxt := context.Background()
	m, ctx := newTreeManager(t)

	addEntity(t, ctx)

	if err := m.LockEntity(ctxt, "entity1"); err != nil {
		t.Fatal(err)
	}

	e, err := ctx.DB.LoadEntity(ctxt, "entity1")
	if err != nil {
		t.Fatal(err)
	}

	if !e.GetMeta().GetLocked() {
		t.Error("Entity not locked")
	}

	if err := m.UnlockEntity(ctxt, "entity1"); err != nil {
		t.Fatal(err)
	}

	e, err = ctx.DB.LoadEntity(ctxt, "entity1")
	if err != nil {
		t.Fatal(err)
	}

	if e.GetMeta().GetLocked() {
		t.Error("Entity not unlocked")
	}
}
