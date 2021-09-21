package interface_test

import (
	"context"
	"testing"

	"github.com/netauth/netauth/internal/db"
)

func TestDeleteGroup(t *testing.T) {
	m, ctx := newTreeManager(t)

	addGroup(t, ctx)

	if err := m.DestroyGroup(context.Background(), "group1"); err != nil {
		t.Fatal(err)
	}

	if _, err := ctx.DB.LoadGroup(context.Background(), "group1"); err != db.ErrUnknownGroup {
		t.Error("Group wasn't deleted")
	}
}
