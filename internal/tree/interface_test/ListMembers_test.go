package interface_test

import (
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/netauth/netauth/internal/db"
)

func TestListMembers(t *testing.T) {
	ctxt := context.Background()
	m, ctx := newTreeManager(t)

	buildSampleTree(t, ctx)
	ctx.DB.(*db.DB).EventUpdateAll()

	// Meta-group ALL, contains all five entities
	mbrs, err := m.ListMembers(ctxt, "ALL")
	if err != nil {
		t.Fatal(err)
	}
	sort.Slice(mbrs, func(i, j int) bool {
		return mbrs[i].GetID() < mbrs[j].GetID()
	})
	for i, e := range mbrs {
		if e.GetID() != fmt.Sprintf("entity%d", i+1) {
			t.Error("Missing an entity")
		}
	}

	// Group 1, should have entity1 and entity2
	mbrs, err = m.ListMembers(ctxt, "group1")
	if err != nil {
		t.Fatal(err)
	}
	sort.Slice(mbrs, func(i, j int) bool {
		return mbrs[i].GetID() < mbrs[j].GetID()
	})
	if len(mbrs) != 2 || mbrs[0].GetID() != "entity1" || mbrs[1].GetID() != "entity2" {
		t.Error("group1 has wrong membership")
	}

	// Group 2, should have entity1 and entity3
	mbrs, err = m.ListMembers(ctxt, "group2")
	if err != nil {
		t.Fatal(err)
	}
	sort.Slice(mbrs, func(i, j int) bool {
		return mbrs[i].GetID() < mbrs[j].GetID()
	})
	if len(mbrs) != 2 || mbrs[0].GetID() != "entity1" || mbrs[1].GetID() != "entity3" {
		t.Error("group2 has wrong membership")
	}

	// Group 4, should have entity1 and NOT entity2
	mbrs, err = m.ListMembers(ctxt, "group4")
	if err != nil {
		t.Fatal(err)
	}
	sort.Slice(mbrs, func(i, j int) bool {
		return mbrs[i].GetID() < mbrs[j].GetID()
	})
	if len(mbrs) != 1 || mbrs[0].GetID() != "entity1" {
		t.Error("group4 has wrong membership")
	}
}
