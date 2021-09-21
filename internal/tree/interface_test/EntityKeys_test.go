package interface_test

import (
	"context"
	"testing"

	"github.com/netauth/netauth/internal/db"
)

func TestManageEntityKeys(t *testing.T) {
	ctxt := context.Background()
	m, ctx := newTreeManager(t)

	addEntity(t, ctx)

	_, err := m.UpdateEntityKeys(ctxt, "entity1", "add", "simple", "secret1")
	if err != nil {
		t.Fatal(err)
	}

	keys, err := m.UpdateEntityKeys(ctxt, "entity1", "list", "", "")
	if err != nil {
		t.Fatal(err)
	}

	if len(keys) != 1 || keys[0] != "SIMPLE:secret1" {
		t.Log(keys)
		t.Error("Key storage error")
	}

	_, err = m.UpdateEntityKeys(ctxt, "entity1", "del", "simple", "")
	if err != nil {
		t.Fatal(err)
	}

	keys, err = m.UpdateEntityKeys(ctxt, "entity1", "list", "", "")
	if err != nil {
		t.Fatal(err)
	}

	if len(keys) != 0 {
		t.Error("Key manipulation error")
	}

}

func TestManageEntityKeysBadEntity(t *testing.T) {
	m, _ := newTreeManager(t)

	_, err := m.UpdateEntityKeys(context.Background(), "entity1", "add", "simple", "secret1")
	if err != db.ErrUnknownEntity {
		t.Fatal(err)
	}
}
