package interface_test

import (
	"testing"

	"github.com/netauth/netauth/internal/db"
)

func TestDeleteGroup(t *testing.T) {
	m, ctx := newTreeManager(t)

	addGroup(t, ctx)

	if err := m.DestroyGroup("group1"); err != nil {
		t.Fatal(err)
	}

	if _, err := ctx.DB.LoadGroup("group1"); err != db.ErrUnknownGroup {
		t.Error("Group wasn't deleted")
	}
}
