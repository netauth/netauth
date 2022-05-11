package interface_test

import (
	"context"
	"testing"

	"github.com/netauth/netauth/internal/db"
)

func TestDeleteEntity(t *testing.T) {
	m, mdb := newTreeManager(t)

	addEntity(t, mdb)

	if err := m.DestroyEntity(context.Background(), "entity1"); err != nil {
		t.Fatal(err)
	}

	if _, err := mdb.LoadEntity(context.Background(), "entity1"); err != db.ErrUnknownEntity {
		t.Error("Entity not deleted")
	}
}
