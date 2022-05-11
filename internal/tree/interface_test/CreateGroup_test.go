package interface_test

import (
	"context"
	"testing"
)

func TestNewGroup(t *testing.T) {
	m, db := newTreeManager(t)

	if err := m.CreateGroup(context.Background(), "group1", "Group 1", "", -1); err != nil {
		t.Fatal(err)
	}

	g, err := db.LoadGroup(context.Background(), "group1")
	if err != nil {
		t.Fatal(err)
	}

	if g.GetDisplayName() != "Group 1" {
		t.Error("Group handling error")
	}
}
