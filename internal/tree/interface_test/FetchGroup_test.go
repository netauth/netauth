package interface_test

import (
	"testing"

	"github.com/golang/protobuf/proto"
)

func TestFetchGroup(t *testing.T) {
	m, ctx := newTreeManager(t)

	addGroup(t, ctx)

	g, err := m.FetchGroup("group1")
	if err != nil {
		t.Fatal(err)
	}

	lg, err := ctx.DB.LoadGroup("group1")
	if err != nil {
		t.Fatal(err)
	}

	if !proto.Equal(g, lg) {
		t.Error("Group handling error")
	}
}
