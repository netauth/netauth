package memdb

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-hclog"

	"github.com/netauth/netauth/internal/db"

	pb "github.com/netauth/protocol"
)

func TestDiscoverEntities(t *testing.T) {
	x, err := New(hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}

	l, err := x.DiscoverEntityIDs()
	if err != nil {
		t.Error(err)
	}

	// At this point there are no entities, so the length should
	// be 0.
	if len(l) != 0 {
		t.Error("DiscoverEntityIDs made up an entity!")
	}

	// We'll save an entity that just has the ID set.  This isn't
	// very realistic, but its the minimum data needed to put a
	// file on disk.
	if err := x.SaveEntity(&pb.Entity{ID: proto.String("foo")}); err != nil {
		t.Error(err)
	}

	// Rerun discovery.
	l, err = x.DiscoverEntityIDs()
	if err != nil {
		t.Error(err)
	}

	// Now there should be one file on disk, and the ID that we've
	// discovered should be 'foo'
	if len(l) != 1 {
		t.Error("DiscoverEntityIDs failed to discover any entities!")
	}
	if l[0] != "foo" {
		t.Errorf("DiscoverEntityIDs discovered the wrong name: '%s'", l[0])
	}
}

func TestSaveLoadDeleteEntity(t *testing.T) {
	x, err := New(hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}

	e := &pb.Entity{ID: proto.String("foo")}

	// Write an entity to disk
	if err := x.SaveEntity(e); err != nil {
		t.Error(err)
	}

	// Load  it back, it  should still be  the same, but  we check
	// this to be sure.
	ne, err := x.LoadEntity("foo")
	if err != nil {
		t.Error(err)
	}
	if !proto.Equal(e, ne) {
		t.Errorf("Loaded entity and original are not equivalent! '%v', '%v'", e, ne)
	}

	// Delete the entity and confirm that loading it returns an
	// error.
	if err := x.DeleteEntity("foo"); err != nil {
		t.Error(err)
	}
	if _, err := x.LoadEntity("foo"); err != db.ErrUnknownEntity {
		t.Error(err)
	}
}

func TestNextEntityNumber(t *testing.T) {
	x, err := New(hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}
	s := []struct {
		ID            string
		number        int32
		nextUIDNumber int32
	}{
		{"foo", 1, 2},
		{"bar", 2, 3},
		{"baz", 65, 66}, // Numbers may be missing in the middle
		{"fuu", 23, 66}, // Later additions shouldn't alter max
	}

	for _, c := range s {
		//  Make sure the entity actually gets added
		e := &pb.Entity{ID: proto.String(c.ID), Number: proto.Int32(c.number)}
		if err := x.SaveEntity(e); err != nil {
			t.Error(err)
		}

		// Validate that after a given mutation the number is
		// still what we expect it to be.
		next, err := x.NextEntityNumber()
		if err != nil {
			t.Error(err)
		}
		if next != c.nextUIDNumber {
			t.Errorf("Wrong next number; got: %v want %v", next, c.nextUIDNumber)
		}
	}
}

func TestSearchEntities(t *testing.T) {
	x, err := New(hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}

	s := []struct {
		ID     string
		secret string
	}{
		{"entity1", "secret1"},
		{"entity2", "secret2"},
		{"entity3", "secret3"},
	}
	for i := range s {
		e := pb.Entity{ID: &s[i].ID, Secret: &s[i].secret}
		if err := x.SaveEntity(&e); err != nil {
			t.Fatal(err)
		}
	}

	_, err = x.SearchEntities(db.SearchRequest{})
	if err != db.ErrBadSearch {
		t.Fatal(err)
	}
	var res []*pb.Entity
	res, err = x.SearchEntities(db.SearchRequest{Expression: "entity1"})
	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 1 || res[0].GetID() != "entity1" {
		t.Log(res)
		t.Error("Result does not match expected singular value")
	}
}

func TestDiscoverGroups(t *testing.T) {
	x, err := New(hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}

	l, err := x.DiscoverGroupNames()
	if err != nil {
		t.Error(err)
	}

	// At this point there are no groups, so the length should
	// be 0.
	if len(l) != 0 {
		t.Error("DiscoverGroupNames made up an entity!")
	}

	// We'll save an entity that just has the ID set.  This isn't
	// very realistic, but its the minimum data needed to put a
	// file on disk.
	if err := x.SaveGroup(&pb.Group{Name: proto.String("foo")}); err != nil {
		t.Error(err)
	}

	// Rerun discovery.
	l, err = x.DiscoverGroupNames()
	if err != nil {
		t.Error(err)
	}

	// Now there should be one file on disk, and the ID that we've
	// discovered should be 'foo'
	if len(l) != 1 {
		t.Error("DiscoverGroupNames failed to discover any groups!")
	}
	if l[0] != "foo" {
		t.Errorf("DiscoverGroupNames discovered the wrong name: '%s'", l[0])
	}
}

func TestGroupSaveLoadDelete(t *testing.T) {
	x, err := New(hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}

	g := &pb.Group{Name: proto.String("foo")}

	// Write an entity to disk
	if err := x.SaveGroup(g); err != nil {
		t.Error(err)
	}

	// Load  it back, it  should still be  the same, but  we check
	// this to be sure.
	ng, err := x.LoadGroup("foo")
	if err != nil {
		t.Error(err)
	}
	if !proto.Equal(g, ng) {
		t.Errorf("Loaded group and original are not equivalent! '%v', '%v'", g, ng)
	}

	// Delete the group and confirm that loading it returns an
	// error.
	if err := x.DeleteGroup("foo"); err != nil {
		t.Error(err)
	}
	if _, err := x.LoadGroup("foo"); err != db.ErrUnknownGroup {
		t.Error(err)
	}
}

func TestDeleteEntityUnknown(t *testing.T) {
	x, err := New(hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}

	if err := x.DeleteEntity("unknown-entity"); err != db.ErrUnknownEntity {
		t.Error(err)
	}
}

func TestDeleteGroupUnknown(t *testing.T) {
	x, err := New(hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}

	if err := x.DeleteGroup("unknown-group"); err != db.ErrUnknownGroup {
		t.Error(err)
	}
}

func TestNextGroupNumber(t *testing.T) {
	x, err := New(hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}
	s := []struct {
		name   string
		number int32
		want   int32
	}{
		{"foo", 1, 2},
		{"bar", 2, 3},
		{"baz", 65, 66}, // Numbers may be missing in the middle
		{"fuu", 23, 66}, // Later additions shouldn't alter max
	}

	for _, c := range s {
		g := &pb.Group{Name: proto.String(c.name), Number: proto.Int32(c.number)}
		if err := x.SaveGroup(g); err != nil {
			t.Error(err)
		}

		// Validate that after a given mutation the number is
		// still what we expect it to be.
		next, err := x.NextGroupNumber()
		if err != nil {
			t.Error(err)
		}
		if next != c.want {
			t.Errorf("Wrong next number; got: %v want %v", next, c.want)
		}
	}
}

func TestSearchGroups(t *testing.T) {
	x, err := New(hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}

	s := []struct {
		Name        string
		DisplayName string
	}{
		{"group1", "One"},
		{"group2", "Two"},
		{"group3", "Three"},
	}
	for i := range s {
		e := pb.Group{Name: &s[i].Name, DisplayName: &s[i].DisplayName}
		if err := x.SaveGroup(&e); err != nil {
			t.Fatal(err)
		}
	}

	_, err = x.SearchGroups(db.SearchRequest{})
	if err != db.ErrBadSearch {
		t.Fatal(err)
	}
	var res []*pb.Group
	res, err = x.SearchGroups(db.SearchRequest{Expression: "group1"})
	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 1 || res[0].GetName() != "group1" {
		t.Log(res)
		t.Error("Result does not match expected singular value")
	}
}

func TestHealthCheck(t *testing.T) {
	x, err := New(hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}

	// Fish out the concrete type to call non-interface methods.
	rx, ok := x.(*MemDB)
	if !ok {
		t.Fatal("Type assertion failed, bad type!")
	}

	if r := rx.healthCheck(); r.OK != true {
		t.Error("hard coded health check somehow changed")
	}
}

func TestLoadSaveEntityErrors(t *testing.T) {
	x, err := New(hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}

	if _, err := x.LoadEntity("load-error"); err != db.ErrInternalError {
		t.Error("Didn't return an error when should have")
	}
	if err := x.SaveEntity(&pb.Entity{ID: proto.String("save-error")}); err != db.ErrInternalError {
		t.Error("Didn't return an error when should have")
	}
}

func TestLoadSaveGroupErrors(t *testing.T) {
	x, err := New(hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}

	if _, err := x.LoadGroup("load-error"); err != db.ErrInternalError {
		t.Error("Didn't return an error when should have")
	}
	if err := x.SaveGroup(&pb.Group{Name: proto.String("save-error")}); err != db.ErrInternalError {
		t.Error("Didn't return an error when should have")
	}
}

// This test case is purely to maintain 100% statement coverage.
func TestCB(t *testing.T) {
	cb()
}
