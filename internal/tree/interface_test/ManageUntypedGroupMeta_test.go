package interface_test

import (
	"testing"

	"github.com/netauth/netauth/internal/db"
)

func TestManageUntypedGroupMeta(t *testing.T) {
	m, ctx := newTreeManager(t)

	addGroup(t, ctx)

	// Add Single Key
	_, err := m.ManageUntypedGroupMeta("group1", "UPSERT", "key1{0}", "value1")
	if err != nil {
		t.Fatal(err)
	}
	uem, err := m.ManageUntypedGroupMeta("group1", "READ", "key1", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(uem) != 1 || uem[0] != "key1{0}:value1" {
		t.Error("Key storage error")
	}

	// Add a second key
	_, err = m.ManageUntypedGroupMeta("group1", "UPSERT", "key1{1}", "value2")
	if err != nil {
		t.Fatal(err)
	}
	uem, err = m.ManageUntypedGroupMeta("group1", "READ", "key1", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(uem) != 2 || uem[1] != "key1{1}:value2" {
		t.Error("Key storage error")
	}

	// Clear the first key
	_, err = m.ManageUntypedGroupMeta("group1", "CLEAREXACT", "key1{0}", "")
	if err != nil {
		t.Fatal(err)
	}
	uem, err = m.ManageUntypedGroupMeta("group1", "READ", "key1", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(uem) != 1 || uem[0] != "key1{1}:value2" {
		t.Error("Key storage error")
	}

	// Clear remaining keys
	_, err = m.ManageUntypedGroupMeta("group1", "CLEARFUZZY", "key1", "")
	if err != nil {
		t.Fatal(err)
	}
	uem, err = m.ManageUntypedGroupMeta("group1", "READ", "key1", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(uem) != 0 {
		t.Error("Key storage error")
	}
}

func TestUntypedGroupMetaBadGroup(t *testing.T) {
	m, _ := newTreeManager(t)

	_, err := m.ManageUntypedGroupMeta("group1", "UPSERT", "key1{0}", "value1")
	if err != db.ErrUnknownGroup {
		t.Fatal(err)
	}
}
