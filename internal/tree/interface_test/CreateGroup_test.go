package interface_test

import (
	"testing"
)

func TestNewGroup(t *testing.T) {
	m, ctx := newTreeManager(t)

	if err := m.NewGroup("group1", "Group 1", "", -1); err != nil {
		t.Fatal(err)
	}

	g, err := ctx.DB.LoadGroup("group1")
	if err != nil {
		t.Fatal(err)
	}

	if g.GetDisplayName() != "Group 1" {
		t.Error("Group handling error")
	}
}
