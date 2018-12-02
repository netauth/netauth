package interface_test

import (
	"testing"
)

func TestSetSecret(t *testing.T) {
	em, ctx := newTreeManager(t)

	addEntity(t, ctx)

	em.SetSecret("entity1", "secret1")

	e, err := ctx.DB.LoadEntity("entity1")
	if err != nil {
		t.Fatal(err)
	}

	if e.GetSecret() != "secret1" {
		t.Error("Secret not set correctly")
	}
}
