package interface_test

import (
	"context"
	"testing"

	"github.com/netauth/netauth/internal/db"
)

func TestManageUntypedEntityMeta(t *testing.T) {
	ctxt := context.Background()
	m, ctx := newTreeManager(t)

	addEntity(t, ctx)

	// Add Single Key
	_, err := m.ManageUntypedEntityMeta(ctxt, "entity1", "UPSERT", "key1{0}", "value1")
	if err != nil {
		t.Fatal(err)
	}
	uem, err := m.ManageUntypedEntityMeta(ctxt, "entity1", "READ", "key1", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(uem) != 1 || uem[0] != "key1{0}:value1" {
		t.Error("Key storage error")
	}

	// Add a second key
	_, err = m.ManageUntypedEntityMeta(ctxt, "entity1", "UPSERT", "key1{1}", "value2")
	if err != nil {
		t.Fatal(err)
	}
	uem, err = m.ManageUntypedEntityMeta(ctxt, "entity1", "READ", "key1", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(uem) != 2 || uem[1] != "key1{1}:value2" {
		t.Error("Key storage error")
	}

	// Clear the first key
	_, err = m.ManageUntypedEntityMeta(ctxt, "entity1", "CLEAREXACT", "key1{0}", "")
	if err != nil {
		t.Fatal(err)
	}
	uem, err = m.ManageUntypedEntityMeta(ctxt, "entity1", "READ", "key1", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(uem) != 1 || uem[0] != "key1{1}:value2" {
		t.Error("Key storage error")
	}

	// Clear remaining keys
	_, err = m.ManageUntypedEntityMeta(ctxt, "entity1", "CLEARFUZZY", "key1", "")
	if err != nil {
		t.Fatal(err)
	}
	uem, err = m.ManageUntypedEntityMeta(ctxt, "entity1", "READ", "key1", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(uem) != 0 {
		t.Error("Key storage error")
	}
}

func TestUntypedEntityMetaBadEntity(t *testing.T) {
	m, _ := newTreeManager(t)

	_, err := m.ManageUntypedEntityMeta(context.Background(), "entity1", "UPSERT", "key1{0}", "value1")
	if err != db.ErrUnknownEntity {
		t.Fatal(err)
	}
}
