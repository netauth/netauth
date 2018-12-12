package util

import (
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/NetAuth/NetAuth/internal/db"

	pb "github.com/NetAuth/Protocol"
)

func TestSearchEntities(t *testing.T) {
	si := NewIndex()

	entities := []pb.Entity{
		pb.Entity{
			ID:     proto.String("entity1"),
			Secret: proto.String("secret"),
			Meta: &pb.EntityMeta{
				GECOS: proto.String("Entity One"),
				Shell: proto.String("/bin/korn"),
			},
		},
		pb.Entity{
			ID:     proto.String("entity2"),
			Secret: proto.String("secret"),
			Meta: &pb.EntityMeta{
				GECOS: proto.String("Entity Two"),
				Shell: proto.String("/bin/fish"),
			},
		},
		pb.Entity{
			ID:     proto.String("entity3"),
			Secret: proto.String("secret"),
			Meta: &pb.EntityMeta{
				GECOS: proto.String("Hazelnut"),
				Shell: proto.String("/bin/korn"),
			},
		},
	}

	for _, e := range entities {
		if err := si.IndexEntity(&e); err != nil {
			t.Fatal(err)
		}
	}

	// Check to make sure secrets didn't get indexed
	r, err := si.SearchEntities(db.SearchRequest{Expression: "secret"})
	if err != nil {
		t.Fatal(err)
	}
	if len(r) > 0 {
		t.Error("Secrets got indexed somehow")
	}

	// Run a test search and make sure there are the right number
	// of answers in it
	r, err = si.SearchEntities(db.SearchRequest{Expression: "meta.Shell:korn"})
	if err != nil {
		t.Fatal(err)
	}
	if len(r) != 2 {
		t.Log(r)
		t.Error("Wrong number of results")
	}
}

func TestSearchGroups(t *testing.T) {
	si := NewIndex()

	groups := []pb.Group{
		pb.Group{
			Name:        proto.String("group1"),
			DisplayName: proto.String("The First Group"),
			UntypedMeta: []string{"UEM:UEM"},
		},
		pb.Group{
			Name:        proto.String("group2"),
			DisplayName: proto.String("The Second Group"),
			UntypedMeta: []string{"UEM:UEM"},
		},
		pb.Group{
			Name:        proto.String("group3"),
			DisplayName: proto.String("This won't match"),
			UntypedMeta: []string{"UEM:UEM"},
		},
	}

	for _, g := range groups {
		if err := si.IndexGroup(&g); err != nil {
			t.Fatal(err)
		}
	}
	// Check to make sure UEM wasn't indexed
	r, err := si.SearchGroups(db.SearchRequest{Expression: "UEM"})
	if err != nil {
		t.Fatal(err)
	}
	if len(r) > 0 {
		t.Fatal("UntypedMeta was indexed")
	}

	// Check a search to make sure its got the right amount of
	// stuff in it.
	r, err = si.SearchGroups(db.SearchRequest{Expression: "DisplayName:Group"})
	if err != nil {
		t.Fatal(err)
	}
	if len(r) != 2 {
		t.Log(r)
		t.Error("result has wrong size")
	}
}
