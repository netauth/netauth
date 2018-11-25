package interface_test

import (
	"testing"

	"github.com/golang/protobuf/proto"
)

func TestGetGroup(t *testing.T) {
	m, ctx := newTreeManager(t)

	addGroup(t, ctx)

	g, err := m.GetGroupByName("group1")
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
