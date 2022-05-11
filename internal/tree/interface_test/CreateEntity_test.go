package interface_test

import (
	"context"
	"testing"
)

func TestNewEntity(t *testing.T) {
	em, db := newTreeManager(t)

	if err := em.CreateEntity(context.Background(), "foo", -1, "foo"); err != nil {
		t.Fatal(err)
	}

	e, err := db.LoadEntity(context.Background(), "foo")
	if err != nil {
		t.Fatal(err)
	}

	if e.GetID() != "foo" {
		t.Error("Entity does not meet saved expectations")
	}
}
