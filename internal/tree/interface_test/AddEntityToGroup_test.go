package interface_test

import (
	"context"
	"testing"
)

func TestAddEntityToGroup(t *testing.T) {
	m, db := newTreeManager(t)

	addEntity(t, db)
	addGroup(t, db)

	if err := m.AddEntityToGroup(context.Background(), "entity1", "group1"); err != nil {
		t.Fatal(err)
	}

	e, err := db.LoadEntity(context.Background(), "entity1")
	if err != nil {
		t.Fatal(err)
	}

	groups := e.GetMeta().GetGroups()
	if len(groups) != 1 || groups[0] != "group1" {
		t.Error("Entity modification error")
	}
}
