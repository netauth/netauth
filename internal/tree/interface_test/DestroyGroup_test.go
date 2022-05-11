package interface_test

import (
	"context"
	"testing"

	"github.com/netauth/netauth/internal/db"
)

func TestDeleteGroup(t *testing.T) {
	m, mdb := newTreeManager(t)

	addGroup(t, mdb)

	if err := m.DestroyGroup(context.Background(), "group1"); err != nil {
		t.Fatal(err)
	}

	if _, err := mdb.LoadGroup(context.Background(), "group1"); err != db.ErrUnknownGroup {
		t.Error("Group wasn't deleted")
	}
}
