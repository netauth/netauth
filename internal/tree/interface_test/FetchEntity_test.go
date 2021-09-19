package interface_test

import (
	"testing"

	"google.golang.org/protobuf/proto"

	"github.com/netauth/netauth/internal/db"
)

func TestFetchEntity(t *testing.T) {
	m, ctx := newTreeManager(t)

	addEntity(t, ctx)

	e, err := m.FetchEntity("entity1")
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

func TestFetchEntityNonExistant(t *testing.T) {
	m, _ := newTreeManager(t)
	if _, err := m.FetchEntity("non-existent"); err != db.ErrUnknownEntity {
		t.Error(err)
	}
}
