package interface_test

import (
	"testing"

	"github.com/netauth/netauth/internal/db"
)

func TestDeleteEntity(t *testing.T) {
	m, ctx := newTreeManager(t)

	addEntity(t, ctx)

	if err := m.DestroyEntity("entity1"); err != nil {
		t.Fatal(err)
	}

	if _, err := ctx.DB.LoadEntity("entity1"); err != db.ErrUnknownEntity {
		t.Error("Entity not deleted")
	}
}
