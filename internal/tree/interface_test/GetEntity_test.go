package interface_test

import (
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/NetAuth/NetAuth/internal/db"
)

func TestGetEntity(t *testing.T) {
	m, ctx := newTreeManager(t)

	addEntity(t, ctx)

	e, err := m.GetEntity("entity1")
	if err != nil {
		t.Fatal(err)
	}

	le, err := ctx.DB.LoadEntity("entity1")
	if err != nil {
		t.Fatal(err)
	}

	le.Secret = proto.String("<REDACTED>")

	if !proto.Equal(e, le) {
		t.Log(e)
		t.Log(le)
		t.Error("Fetched entity is not equivalent")
	}
}

func TestGetEntityNonExistant(t *testing.T) {
	m, _ := newTreeManager(t)
	if _, err := m.GetEntity("non-existant"); err != db.ErrUnknownEntity {
		t.Error(err)
	}
}
