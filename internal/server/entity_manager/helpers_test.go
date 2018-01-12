package entity_manager

import (
	"testing"

	"github.com/golang/protobuf/proto"
)

func TestSafeCopyEntity(t *testing.T) {
	em := New()

	if err := em.newEntity("foo", -1, "bar"); err != nil {
		t.Error(err)
	}

	e, err := em.getEntityByID("foo")
	if err != nil {
		t.Error(err)
	}

	ne, err := safeCopyEntity(e)
	if err != nil {
		t.Error(err)
	}

	// The normal way to do this would be to check if the proto is
	// the same, but here we need to check if two fields are
	// different, then make sure that everything else is the same.
	if e.GetSecret() == ne.GetSecret() {
		t.Error("Secret field not obscured!")
	}

	e.Secret = proto.String("")
	ne.Secret = proto.String("")

	if !proto.Equal(e, ne) {
		t.Error("Entity values not otherwise equal!")
	}
}
