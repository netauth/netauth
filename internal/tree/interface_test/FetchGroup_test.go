package interface_test

import (
	"context"
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestFetchGroup(t *testing.T) {
	ctxt := context.Background()
	m, ctx := newTreeManager(t)

	addGroup(t, ctx)

	g, err := m.FetchGroup(ctxt, "group1")
	if err != nil {
		t.Fatal(err)
	}

	lg, err := ctx.DB.LoadGroup(ctxt, "group1")
	if err != nil {
		t.Fatal(err)
	}

	if !proto.Equal(g, lg) {
		t.Error("Group handling error")
	}
}
