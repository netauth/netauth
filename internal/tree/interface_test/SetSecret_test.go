package interface_test

import (
	"context"
	"testing"
)

func TestSetSecret(t *testing.T) {
	ctxt := context.Background()
	em, mdb := newTreeManager(t)

	addEntity(t, mdb)

	em.SetSecret(ctxt, "entity1", "secret1")

	e, err := mdb.LoadEntity(ctxt, "entity1")
	if err != nil {
		t.Fatal(err)
	}

	if e.GetSecret() != "secret1" {
		t.Error("Secret not set correctly")
	}
}
