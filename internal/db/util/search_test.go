package util

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-hclog"

	"github.com/netauth/netauth/internal/db"

	pb "github.com/netauth/protocol"
)

func dummyEntityLoader(e string) (*pb.Entity, error) {
	switch e {
	case "entity1":
		return &pb.Entity{
			ID:     proto.String("entity1"),
			Secret: proto.String("secret"),
			Meta: &pb.EntityMeta{
				GECOS: proto.String("Entity One"),
				Shell: proto.String("/bin/korn"),
			},
		}, nil
	default:
		return nil, db.ErrUnknownEntity
	}
}

func dummyGroupLoader(g string) (*pb.Group, error) {
	switch g {
	case "group1":
		return &pb.Group{
			Name:        proto.String("group1"),
			DisplayName: proto.String("The First Group"),
			UntypedMeta: []string{"UEM:UEM"},
		}, nil
	default:
		return nil, db.ErrUnknownGroup
	}
}

// This test case is explicitly for checking the call path in the
// coverage test output.  The intent is to validate that the loader
// will bail out if the callback is not configured.  If the callback
// has not been configured, it will return in the correct case, or
// panic in the incorrect case.
func TestIndexCallbackUnconfigured(t *testing.T) {
	si := NewIndex(hclog.NewNullLogger())
	si.IndexCallback(db.Event{Type: db.EventEntityCreate, PK: "entity1"})
}

func TestIndexCallbackEntity(t *testing.T) {
	si := NewIndex(hclog.NewNullLogger())
	si.ConfigureCallback(dummyEntityLoader, dummyGroupLoader)

	// Check that the entity isn't present
	r, err := si.SearchEntities(db.SearchRequest{Expression: "ID:entity1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(r) != 0 {
		t.Error("Entity exists in a nil SearchIndex")
	}

	// Index the entity
	si.IndexCallback(db.Event{Type: db.EventEntityCreate, PK: "entity1"})

	// Check for entity being present in results
	r, err = si.SearchEntities(db.SearchRequest{Expression: "ID:entity1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(r) != 1 {
		t.Error("Exact match wasn't returned")
	}

	// Index an entity that doesn't exist.  This is primarily to
	// make sure that the loader doesn't explode when trying to
	// fetch an entity that doesn't exist.
	si.IndexCallback(db.Event{Type: db.EventEntityCreate, PK: "entity2"})

	// Fire a "delete" and make sure the entity drops out of the search results
	si.IndexCallback(db.Event{Type: db.EventEntityDestroy, PK: "entity1"})
	r, err = si.SearchEntities(db.SearchRequest{Expression: "ID:entity1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(r) != 0 {
		t.Error("Entity exists in a nil SearchIndex")
	}
}

func TestIndexCallbackGroup(t *testing.T) {
	si := NewIndex(hclog.NewNullLogger())
	si.ConfigureCallback(dummyEntityLoader, dummyGroupLoader)

	// Check that the group isn't present
	r, err := si.SearchGroups(db.SearchRequest{Expression: "Name:group1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(r) != 0 {
		t.Error("Group exists in a nil SearchIndex")
	}

	// Index the group
	si.IndexCallback(db.Event{Type: db.EventGroupCreate, PK: "group1"})

	// Check for group being present in results
	r, err = si.SearchGroups(db.SearchRequest{Expression: "Name:group1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(r) != 1 {
		t.Error("Exact match wasn't returned")
	}

	// Index an group that doesn't exist.  This is primarily to
	// make sure that the loader doesn't explode when trying to
	// fetch an group that doesn't exist.
	si.IndexCallback(db.Event{Type: db.EventGroupCreate, PK: "group2"})

	// Fire a "delete" and make sure the group drops out of the search results
	si.IndexCallback(db.Event{Type: db.EventGroupDestroy, PK: "group1"})
	r, err = si.SearchGroups(db.SearchRequest{Expression: "Name:group1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(r) != 0 {
		t.Error("Group exists in a nil SearchIndex")
	}
}

func TestSearchEntities(t *testing.T) {
	si := NewIndex(hclog.NewNullLogger())

	entities := []pb.Entity{
		{
			ID:     proto.String("entity1"),
			Secret: proto.String("secret"),
			Meta: &pb.EntityMeta{
				GECOS: proto.String("Entity One"),
				Shell: proto.String("/bin/korn"),
			},
		},
		{
			ID:     proto.String("entity2"),
			Secret: proto.String("secret"),
			Meta: &pb.EntityMeta{
				GECOS: proto.String("Entity Two"),
				Shell: proto.String("/bin/fish"),
			},
		},
		{
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

	// Remove entity2 from the index and search for fish, should
	// yield no results
	if err := si.DeleteEntity(&entities[1]); err != nil {
		t.Error(err)
	}
	r, err = si.SearchEntities(db.SearchRequest{Expression: "meta.Shell:fish"})
	if err != nil {
		t.Fatal(err)
	}
	if len(r) != 0 {
		t.Log(r)
		t.Error("Got results back for an entity which was removed")
	}

}

func TestSearchEntitiesBadRequest(t *testing.T) {
	si := NewIndex(hclog.NewNullLogger())

	r, err := si.SearchEntities(db.SearchRequest{})
	if err != db.ErrBadSearch || r != nil {
		t.Error(err)
	}
}

func TestSearchGroups(t *testing.T) {
	si := NewIndex(hclog.NewNullLogger())

	groups := []pb.Group{
		{
			Name:        proto.String("group1"),
			DisplayName: proto.String("The First Group"),
			UntypedMeta: []string{"UEM:UEM"},
		},
		{
			Name:        proto.String("group2"),
			DisplayName: proto.String("The Second Group"),
			UntypedMeta: []string{"UEM:UEM"},
		},
		{
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

	// Remove group3 and check that its missing when searching for
	// 'match'
	if err := si.DeleteGroup(&groups[2]); err != nil {
		t.Fatal(err)
	}
	r, err = si.SearchGroups(db.SearchRequest{Expression: "DisplayName:match"})
	if err != nil {
		t.Fatal(err)
	}
	if len(r) != 0 {
		t.Log(r)
		t.Error("Search results contained deleted group")
	}
}

func TestSearchGroupsBadRequest(t *testing.T) {
	si := NewIndex(hclog.NewNullLogger())

	r, err := si.SearchGroups(db.SearchRequest{})
	if err != db.ErrBadSearch || r != nil {
		t.Error(err)
	}
}

func TestExtractDocIDsNullResult(t *testing.T) {
	if res := extractDocIDs(nil); res != nil {
		t.Error("Got a non-nil response from a nil result")
	}
}
