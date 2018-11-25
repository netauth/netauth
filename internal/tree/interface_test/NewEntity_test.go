package interface_test

import (
	"testing"

)

func TestNewEntity(t *testing.T) {
	em, ctx := newTreeManager(t)

	if err := em.NewEntity("foo", -1, "foo"); err != nil {
		t.Fatal(err)
	}

	e, err := ctx.DB.LoadEntity("foo")
	if err != nil {
		t.Fatal(err)
	}

	if e.GetID() != "foo" {
		t.Error("Entity does not meet saved expectations")
	}
}
